package raydium

import (
	"fmt"
	"indexer/src/indexer"
	"indexer/src/programs"
	"indexer/src/util/solana/transactions"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

func PopulateInitialize(reader *transactions.Reader, flatIndex uint8) (s *Swap, err error) {
	fmt.Println("initialize", reader.GetSignature())
	panic("not implemented")
}

func PopulateInitialize2(reader *transactions.Reader, flatIndex uint8) (i *Initialize2, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 21) {
		return nil, fmt.Errorf("accounts length for MonitorStep must be either 17 or 18")
	}

	i = &Initialize2{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	i.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[4])
	if err != nil {
		return nil, err
	}
	i.Minter, err = reader.GetAccountAtIndex(inst.Accounts[17])
	if err != nil {
		return nil, err
	}
	i.CoinMint, err = reader.GetAccountAtIndex(inst.Accounts[8])
	if err != nil {
		return nil, err
	}
	i.PoolCoinTokenAccount, err = reader.GetAccountAtIndex(inst.Accounts[10])
	if err != nil {
		return nil, err
	}
	i.PcMint, err = reader.GetAccountAtIndex(inst.Accounts[9])
	if err != nil {
		return nil, err
	}
	i.PoolPcTokenAccount, err = reader.GetAccountAtIndex(inst.Accounts[11])
	if err != nil {
		return nil, err
	}
	i.LpMint, err = reader.GetAccountAtIndex(inst.Accounts[7])
	if err != nil {
		return nil, err
	}

	initialize2Inst := &InstructionInitialize2{}
	err = bin.NewBinDecoder(inst.Data).Decode(initialize2Inst)
	if err != nil {
		return nil, err
	}

	i.Nonce = initialize2Inst.Nonce
	i.InitPcAmount = initialize2Inst.InitPcAmount
	i.InitCoinAmount = initialize2Inst.InitCoinAmount

	userLpTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[20])
	if err != nil {
		return nil, err
	}

	lpMintInst, err := reader.FindInstruction(flatIndex+1, func(i2 solana.CompiledInstruction, idx uint8) bool {
		if len(i2.Data) == 0 {
			return false
		}
		accData, _ := reader.GetAccountAtIndex(i2.ProgramIDIndex)
		if !accData.Equals(solana.TokenProgramID) {
			return false
		}

		discriminator, err := programs.GetInstructionDiscriminator(i2.Data)
		if err != nil {
			return false
		}

		return discriminator == 7
	})
	if err != nil {
		return nil, err
	}

	accounts, err := reader.GetAccountsAtIndices(lpMintInst.Accounts)
	if err != nil {
		return nil, err
	}

	accountMeta := make([]*solana.AccountMeta, 0, len(accounts))
	for _, acc := range accounts {
		accountMeta = append(accountMeta, &solana.AccountMeta{PublicKey: acc})
	}

	parsed, err := token.DecodeInstruction(accountMeta, lpMintInst.Data)
	if err != nil {
		return nil, err
	}

	mintTo, ok := parsed.Impl.(*token.MintTo)
	if !ok {
		return nil, fmt.Errorf("error casting insttruction to mint2")
	}

	dest, err := reader.GetAccountAtIndex(lpMintInst.Accounts[1])
	if err != nil {
		return nil, err
	}

	if !dest.Equals(userLpTokenAccount) {
		return nil, fmt.Errorf("mintTo instruction does not mint to user lp token account")
	}
	i.LpAmount = *mintTo.Amount

	return i, nil
}

func PopulateMonitorStep(reader *transactions.Reader, flatIndex uint8) (m *MonitorStep, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 18) {
		return nil, fmt.Errorf("accounts length for MonitorStep must be either 17 or 18")
	}

	m = &MonitorStep{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	m.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[3])
	if err != nil {
		return nil, err
	}

	monitorStepInst := &InstructionMonitorStep{}
	err = bin.NewBinDecoder(inst.Data).Decode(monitorStepInst)
	if err != nil {
		return nil, err
	}
	m.PlanOrderLimit = monitorStepInst.PlanOrderLimit
	m.PlaceOrderLimit = monitorStepInst.PlaceOrderLimit
	m.CancelOrderLimit = monitorStepInst.CancelOrderLimit

	return m, nil
}

func PopulateAddLiquidity(reader *transactions.Reader, flatIndex uint8) (a *AddLiquidity, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 14) {
		return nil, fmt.Errorf("accounts length for AddLiquidity must be either 17 or 18")
	}

	a = &AddLiquidity{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	a.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	a.Minter, err = reader.GetAccountAtIndex(inst.Accounts[12])
	if err != nil {
		return nil, err
	}

	addLiquidityInst := &InstructionAddLiquidity{}
	err = bin.NewBinDecoder(inst.Data).Decode(addLiquidityInst)
	if err != nil {
		return nil, err
	}

	poolCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}
	poolPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[7])
	if err != nil {
		return nil, err
	}
	userCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[9])
	if err != nil {
		return nil, err
	}
	userPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[10])
	if err != nil {
		return nil, err
	}

	poolCoinTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Destination.Equals(poolCoinTokenAccount) && t.Source.Equals(userCoinTokenAccount) && t.FlattenedInstructionIndex > flatIndex
	})
	if err != nil {
		return nil, err
	}
	poolPcTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Destination.Equals(poolPcTokenAccount) && t.Source.Equals(userPcTokenAccount) && t.FlattenedInstructionIndex > flatIndex
	})
	if err != nil {
		return nil, err
	}

	if poolCoinTransfer.FlattenedInstructionIndex+1 != poolPcTransfer.FlattenedInstructionIndex {
		return nil, fmt.Errorf("payment transaction ordering invalid")
	}

	accounts, err := reader.GetAccountsAtIndices(inst.Accounts)
	if err != nil {
		return nil, err
	}

	accountMeta := make([]*solana.AccountMeta, 0, len(accounts))
	for _, acc := range accounts {
		accountMeta = append(accountMeta, &solana.AccountMeta{PublicKey: acc})
	}

	expectedMintInstructionIndexFlat := poolPcTransfer.FlattenedInstructionIndex + 1
	if int(expectedMintInstructionIndexFlat) >= len(reader.GetFlattenedInstructions()) {
		return nil, fmt.Errorf("expected mint instruction index out of bounds")
	}

	nextInst, err := reader.GetInstructionAtFlattenedIndex(expectedMintInstructionIndexFlat)
	if err != nil {
		return nil, err
	}
	mintInstruction, err := token.DecodeInstruction(accountMeta, nextInst.Data)
	if err != nil {
		return nil, err
	}

	mintToInstr, ok := mintInstruction.Impl.(*token.MintTo)
	if !ok {
		return nil, fmt.Errorf("failed to cast to mintTo, instead recieved type %T", mintInstruction.Impl)
	}

	a.LpTokenAmount = *mintToInstr.Amount

	return a, nil
}

func PopulateRemoveLiquidity(reader *transactions.Reader, flatIndex uint8) (r *RemoveLiquidity, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 22) {
		return nil, fmt.Errorf("accounts length for RemoveLiquidity must be either 17 or 18")
	}

	r = &RemoveLiquidity{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	r.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[0])
	if err != nil {
		return nil, err
	}
	r.Owner, err = reader.GetAccountAtIndex(inst.Accounts[18])
	if err != nil {
		return nil, err
	}

	removeLiquidityInst := &InstructionRemoveLiquidity{}
	err = bin.NewBinDecoder(inst.Data).Decode(removeLiquidityInst)
	if err != nil {
		return nil, err
	}
	r.LpTokenAmount = removeLiquidityInst.Amount

	poolCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}
	userCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[16])
	if err != nil {
		return nil, err
	}
	coinTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Source.Equals(poolCoinTokenAccount) && t.Destination.Equals(userCoinTokenAccount)
	})
	if err != nil {
		return nil, err
	}
	r.AmountBase = coinTransfer.Amount

	poolPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[7])
	if err != nil {
		return nil, err
	}
	userPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[17])
	if err != nil {
		return nil, err
	}
	pcTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Source.Equals(poolPcTokenAccount) && t.Destination.Equals(userPcTokenAccount)
	})
	if err != nil {
		return nil, err
	}
	r.AmountQuote = pcTransfer.Amount

	return r, nil
}

func PopulateWithdrawPnl(reader *transactions.Reader, flatIndex uint8) (w *WithdrawPnl, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 17) {
		return nil, fmt.Errorf("accounts length for WithdrawPnl must be either 17 or 18")
	}

	w = &WithdrawPnl{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	w.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	w.Owner, err = reader.GetAccountAtIndex(inst.Accounts[9])
	if err != nil {
		return nil, err
	}

	poolCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[5])
	if err != nil {
		return nil, err
	}
	poolPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}
	userCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[7])
	if err != nil {
		return nil, err
	}
	userPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[8])
	if err != nil {
		return nil, err
	}

	poolCoinTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Source.Equals(poolCoinTokenAccount) && t.Destination.Equals(userCoinTokenAccount) && t.FlattenedInstructionIndex > flatIndex
	})
	if err != nil {
		return nil, err
	}
	poolPcTransfer, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return t.Source.Equals(poolPcTokenAccount) && t.Destination.Equals(userPcTokenAccount) && t.FlattenedInstructionIndex > flatIndex
	})
	if err != nil {
		return nil, err
	}

	if poolCoinTransfer.FlattenedInstructionIndex+1 != poolPcTransfer.FlattenedInstructionIndex {
		return nil, fmt.Errorf("payment transaction ordering invalid")
	}

	w.AmountBase = poolCoinTransfer.Amount
	w.AmountQuote = poolPcTransfer.Amount

	return w, nil
}

func PopulateSwapExactAmountIn(reader *transactions.Reader, flatIndex uint8) (s *Swap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 17 || accountsLen == 18) {
		return nil, fmt.Errorf("accounts length for SwapExactAmountIn must be either 17 or 18")
	}

	s = &Swap{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}

	offset := accountsLen - 17
	s.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	s.Maker, err = reader.GetAccountAtIndex(inst.Accounts[16+offset])
	if err != nil {
		return nil, err
	}

	poolCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[4+offset])
	if err != nil {
		return nil, err
	}
	poolPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[5+offset])
	if err != nil {
		return nil, err
	}
	userSourceTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[14+offset])
	if err != nil {
		return nil, err
	}
	userDestinationTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[15+offset])
	if err != nil {
		return nil, err
	}

	swapExactAmountInParams := &InstructionSwapExactAmountIn{}
	err = bin.NewBinDecoder(inst.Data).Decode(swapExactAmountInParams)
	if err != nil {
		return nil, err
	}

	// Need to find transfer from `userSourceTokenAccount` to either `poolCoinTokenAccount` or `poolPcTokenAccount`.
	toPoolFromUserSource, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return (t.FlattenedInstructionIndex > flatIndex) && t.Source.Equals(userSourceTokenAccount) && (t.Destination.Equals(poolCoinTokenAccount) || t.Destination.Equals(poolPcTokenAccount)) && t.Amount == swapExactAmountInParams.AmountIn
	})
	if err != nil {
		return nil, err
	}

	fromPoolToUserDest, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return (t.FlattenedInstructionIndex > flatIndex) && t.Destination.Equals(userDestinationTokenAccount) && (t.Source.Equals(poolCoinTokenAccount) || t.Source.Equals(poolPcTokenAccount))
	})
	if err != nil {
		return nil, err
	}

	if toPoolFromUserSource.FlattenedInstructionIndex+1 != fromPoolToUserDest.FlattenedInstructionIndex {
		return nil, fmt.Errorf("payment transaction ordering invalid")
	}

	isSell := toPoolFromUserSource.Destination.Equals(poolCoinTokenAccount)

	if isSell {
		s.AmountBase = int64(toPoolFromUserSource.Amount)
		s.AmountQuote = -int64(fromPoolToUserDest.Amount)
	} else {
		s.AmountBase = -int64(toPoolFromUserSource.Amount)
		s.AmountQuote = int64(fromPoolToUserDest.Amount)
	}

	return s, nil
}

func PopulateSwapExactAmountOut(reader *transactions.Reader, flatIndex uint8) (s *Swap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	accountsLen := len(inst.Accounts)
	if !(accountsLen == 17 || accountsLen == 18) {
		return nil, fmt.Errorf("accounts length for SwapExactAmountOut must be either 17 or 18")
	}

	s = &Swap{
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}

	offset := accountsLen - 17
	s.PoolIdentifier, err = reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	s.Maker, err = reader.GetAccountAtIndex(inst.Accounts[16+offset])
	if err != nil {
		return nil, err
	}

	poolCoinTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[4+offset])
	if err != nil {
		return nil, err
	}
	poolPcTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[5+offset])
	if err != nil {
		return nil, err
	}
	userSourceTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[14+offset])
	if err != nil {
		return nil, err
	}
	userDestinationTokenAccount, err := reader.GetAccountAtIndex(inst.Accounts[15+offset])
	if err != nil {
		return nil, err
	}

	swapExactOutParams := &InstructionSwapExactAmountOut{}
	err = bin.NewBinDecoder(inst.Data).Decode(swapExactOutParams)
	if err != nil {
		return nil, err
	}

	// Need to find transfer from `userSourceTokenAccount` to either `poolCoinTokenAccount` or `poolPcTokenAccount`.
	toPoolFromUserSource, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return (t.FlattenedInstructionIndex > flatIndex) && t.Source.Equals(userSourceTokenAccount) && (t.Destination.Equals(poolCoinTokenAccount) || t.Destination.Equals(poolPcTokenAccount))
	})
	if err != nil {
		return nil, err
	}

	fromPoolToUserDest, err := reader.FindTransfer(func(t *transactions.TokenTransfer) bool {
		return (t.FlattenedInstructionIndex > flatIndex) && t.Destination.Equals(userDestinationTokenAccount) && (t.Source.Equals(poolCoinTokenAccount) || t.Source.Equals(poolPcTokenAccount)) && t.Amount == swapExactOutParams.AmountOut
	})
	if err != nil {
		return nil, err
	}

	if toPoolFromUserSource.FlattenedInstructionIndex+1 != fromPoolToUserDest.FlattenedInstructionIndex {
		return nil, fmt.Errorf("payment transaction ordering invalid")
	}

	isSell := toPoolFromUserSource.Destination.Equals(poolCoinTokenAccount)

	if isSell {
		s.AmountBase = int64(toPoolFromUserSource.Amount)
		s.AmountQuote = -int64(fromPoolToUserDest.Amount)
	} else {
		s.AmountBase = -int64(toPoolFromUserSource.Amount)
		s.AmountQuote = int64(fromPoolToUserDest.Amount)
	}

	return s, nil
}

func populateMetadata(reader *transactions.Reader, flatIndex uint8) *indexer.InstructionMetadata {
	return &indexer.InstructionMetadata{
		Slot:             reader.GetSlot(),
		TransactionIndex: reader.GetTransactionIndex(),
		InstructionIndex: flatIndex,
		Signature:        reader.GetSignature(),
		Timestamp:        reader.GetTimestamp(),
	}
}
