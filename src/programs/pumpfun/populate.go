package pumpfun

import (
	"fmt"
	"indexer/src/indexer"
	"indexer/src/util/solana/transactions"

	bin "github.com/gagliardetto/binary"
)

func PopulateCreate(reader *transactions.Reader, flatIndex uint8) (c *Create, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 12 && len(inst.Accounts) != 14 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun create transaction for tx: %s", reader.GetSignature().String())
	}
	createInst := &InstructionCreate{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(createInst)
	if err != nil {
		return nil, err
	}
	c = &Create{
		Name:                string(createInst.Name),
		Symbol:              string(createInst.Symbol),
		MetadataURI:         string(createInst.URI),
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	c.Mint, err = reader.GetAccountAtIndex(inst.Accounts[0])
	if err != nil {
		return nil, err
	}
	c.BondingCurve, err = reader.GetAccountAtIndex(inst.Accounts[2])
	if err != nil {
		return nil, err
	}
	c.AssociatedBondingCurve, err = reader.GetAccountAtIndex(inst.Accounts[3])
	if err != nil {
		return nil, err
	}
	c.Deployer, err = reader.GetAccountAtIndex(inst.Accounts[7])
	if err != nil {
		return nil, err
	}
	c.MetadataSlot, err = reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}
	return c, nil
}

func PopulateBuy(reader *transactions.Reader, flatIndex uint8) (b *Swap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 10 && len(inst.Accounts) != 12 && len(inst.Accounts) != 14 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun buy transaction for tx: %s", reader.GetSignature().String())
	}
	buyInst := &InstructionBuy{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(buyInst)
	if err != nil {
		return nil, err
	}
	b = &Swap{
		TokenAmount:         -int64(buyInst.Amount),
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	b.Mint, err = reader.GetAccountAtIndex(inst.Accounts[2])
	if err != nil {
		return nil, err
	}
	bondingCurve, err := reader.GetAccountAtIndex(inst.Accounts[3])
	if err != nil {
		return nil, err
	}
	b.MakerTokenAccount, err = reader.GetAccountAtIndex(inst.Accounts[5])
	if err != nil {
		return nil, err
	}
	b.Maker, err = reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}

	tokenPayment, err := reader.FindSolTransfer(func(t *transactions.SolTransfer) bool {
		return t.FlattenedInstructionIndex > flatIndex && t.Source.Equals(b.Maker) && t.Destination.Equals(bondingCurve)
	})
	if err != nil {
		fmt.Println("payment", reader.GetSignature())
		return nil, err
	}
	b.SolAmount = int64(tokenPayment.Amount)
	feeAccount, err := reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	feePayment, err := reader.FindSolTransfer(func(t *transactions.SolTransfer) bool {
		return t.FlattenedInstructionIndex > tokenPayment.FlattenedInstructionIndex && t.Source.Equals(b.Maker) && t.Destination.Equals(feeAccount)
	})
	if err != nil {
		fmt.Println("fee", reader.GetSignature())
		return nil, err
	}
	b.Fee = feePayment.Amount
	return b, nil
}

func PopulateSell(reader *transactions.Reader, flatIndex uint8) (s *Swap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 10 && len(inst.Accounts) != 12 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun sell transaction for tx: %s", reader.GetSignature().String())
	}
	sellInst := &InstructionSell{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(sellInst)
	if err != nil {
		return nil, err
	}
	s = &Swap{
		TokenAmount:         int64(sellInst.Amount),
		InstructionMetadata: populateMetadata(reader, flatIndex),
	}
	s.Mint, err = reader.GetAccountAtIndex(inst.Accounts[2])
	if err != nil {
		return nil, err
	}
	s.MakerTokenAccount, err = reader.GetAccountAtIndex(inst.Accounts[5])
	if err != nil {
		return nil, err
	}
	s.Maker, err = reader.GetAccountAtIndex(inst.Accounts[6])
	if err != nil {
		return nil, err
	}
	s.SolAmount, err = reader.GetSolBalanceDelta(inst.Accounts[3])
	if err != nil {
		return nil, err
	}
	if s.SolAmount > 0 {
		return nil, fmt.Errorf("positive bonding curve balance delta invalid for tx: %s", reader.GetSignature().String())
	}
	fee, err := reader.GetSolBalanceDelta(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	if fee < 0 {
		return nil, fmt.Errorf("invalid fee account balance delta")
	}
	s.Fee = uint64(fee)
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
