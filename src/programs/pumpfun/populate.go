package pumpfun

import (
	"fmt"
	"indexer/src/db/db_types"
	"indexer/src/util/solana/transactions"

	bin "github.com/gagliardetto/binary"
)

func PopulateCreate(reader *transactions.Reader, flatIndex uint8) (c *db_types.PumpFunCreate, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 12 && len(inst.Accounts) != 14 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun create transaction for tx: %s", reader.GetSignature().String())
	}
	createInst := &InstructionCreate{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(createInst)
	if err != nil {
		return nil, err
	}
	c = &db_types.PumpFunCreate{
		Name:                string(createInst.Name),
		Symbol:              string(createInst.Symbol),
		MetadataURI:         string(createInst.URI),
		InstructionMetadata: db_types.PopulateMetadata(reader, flatIndex),
	}
	c.Mint, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[0]))
	if err != nil {
		return nil, err
	}
	c.BondingCurve, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[2]))
	if err != nil {
		return nil, err
	}
	c.AssociatedBondingCurve, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[3]))
	if err != nil {
		return nil, err
	}
	c.Deployer, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[7]))
	if err != nil {
		return nil, err
	}
	c.MetadataSlot, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[6]))
	if err != nil {
		return nil, err
	}
	return c, nil
}

func PopulateBuy(reader *transactions.Reader, flatIndex uint8) (b *db_types.PumpFunSwap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 10 && len(inst.Accounts) != 12 && len(inst.Accounts) != 14 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun buy transaction for tx: %s", reader.GetSignature().String())
	}
	buyInst := &InstructionBuy{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(buyInst)
	if err != nil {
		return nil, err
	}
	b = &db_types.PumpFunSwap{
		TokenAmount:         -int64(buyInst.Amount),
		InstructionMetadata: db_types.PopulateMetadata(reader, flatIndex),
	}
	b.Mint, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[2]))
	if err != nil {
		return nil, err
	}
	bondingCurve, err := reader.GetAccountAtIndex(inst.Accounts[3])
	if err != nil {
		return nil, err
	}
	b.MakerTokenAccount, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[5]))
	if err != nil {
		return nil, err
	}
	b.Maker, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[6]))
	if err != nil {
		return nil, err
	}

	tokenPayment, err := reader.FindSolTransfer(func(t *transactions.SolTransfer) bool {
		return t.FlattenedInstructionIndex > flatIndex && t.Source.Equals(b.Maker.PublicKey()) && t.Destination.Equals(bondingCurve)
	})
	if err != nil {
		return nil, err
	}
	b.SolAmount = int64(tokenPayment.Amount)
	feeAccount, err := reader.GetAccountAtIndex(inst.Accounts[1])
	if err != nil {
		return nil, err
	}
	feePayment, err := reader.FindSolTransfer(func(t *transactions.SolTransfer) bool {
		return t.FlattenedInstructionIndex > tokenPayment.FlattenedInstructionIndex && t.Source.Equals(b.Maker.PublicKey()) && t.Destination.Equals(feeAccount)
	})
	if err != nil {
		return nil, err
	}
	b.Fee = feePayment.Amount
	return b, nil
}

func PopulateSell(reader *transactions.Reader, flatIndex uint8) (s *db_types.PumpFunSwap, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	if len(inst.Accounts) != 10 && len(inst.Accounts) != 12 {
		return nil, fmt.Errorf("invalid accounts length for pump.fun sell transaction for tx: %s", reader.GetSignature().String())
	}
	sellInst := &InstructionSell{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(sellInst)
	if err != nil {
		return nil, err
	}
	s = &db_types.PumpFunSwap{
		TokenAmount:         int64(sellInst.Amount),
		InstructionMetadata: db_types.PopulateMetadata(reader, flatIndex),
	}
	s.Mint, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[2]))
	if err != nil {
		return nil, err
	}
	s.MakerTokenAccount, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[5]))
	if err != nil {
		return nil, err
	}
	s.Maker, err = db_types.ToPublicKeyErr(reader.GetAccountAtIndex(inst.Accounts[6]))
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

func PopulateSetParams(reader *transactions.Reader, flatIndex uint8) (s *db_types.PumpFunSetParams, err error) {
	inst, _ := reader.GetInstructionAtFlattenedIndex(flatIndex)
	setParamsInst := &InstructionSetParams{}
	err = bin.NewDecoderWithEncoding(inst.Data, bin.EncodingBorsh).Decode(setParamsInst)
	if err != nil {
		return nil, err
	}
	return &db_types.PumpFunSetParams{
		FeeRecipient:                db_types.PublicKey(setParamsInst.FeeRecipient),
		InitialVirtualTokenReserves: setParamsInst.InitialVirtualTokenReserves,
		InitialVirtualSolReserves:   setParamsInst.InitialVirtualSolReserves,
		InitialRealTokenReserves:    setParamsInst.InitialRealTokenReserves,
		TokenTotalSupply:            setParamsInst.TokenTotalSupply,
		FeeBasisPoints:              setParamsInst.FeeBasisPoints,
		InstructionMetadata:         db_types.PopulateMetadata(reader, flatIndex),
	}, nil
}
