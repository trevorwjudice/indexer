package transactions

import (
	"errors"
	"fmt"
	"indexer/src/util/solana/logparse"
	"strconv"
	"time"

	bin "github.com/gagliardetto/binary"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

type Reader struct {
	block   *rpc.GetBlockResult
	txIndex int
	slot    uint64

	tx                    *rpc.TransactionWithMeta
	parsed                *solana.Transaction
	accountKeys           []solana.PublicKey
	tokenDeltas           *AccountBalanceDeltaMap
	tokenTransfers        []*TokenTransfer
	solTransfers          []*SolTransfer
	flattenedInstructions []solana.CompiledInstruction
	accounts              map[solana.PublicKey]struct{}
}

func NewReader(block *rpc.GetBlockResult, slot uint64, transactionIndex int) (*Reader, error) {
	if transactionIndex >= len(block.Transactions) {
		return nil, fmt.Errorf("transaction index out of range")
	}

	tx := &block.Transactions[transactionIndex]
	parsed, err := tx.GetTransaction()
	if err != nil {
		return nil, err
	}
	res := &Reader{
		block:       block,
		tx:          tx,
		slot:        slot,
		txIndex:     transactionIndex,
		parsed:      parsed,
		accountKeys: make([]solana.PublicKey, 0, len(parsed.Message.AccountKeys)+len(tx.Meta.LoadedAddresses.ReadOnly)+len(tx.Meta.LoadedAddresses.Writable)),
		accounts:    make(map[solana.PublicKey]struct{}),
	}

	res.accountKeys = append(res.accountKeys, parsed.Message.AccountKeys...)
	res.accountKeys = append(res.accountKeys, tx.Meta.LoadedAddresses.Writable...)
	res.accountKeys = append(res.accountKeys, tx.Meta.LoadedAddresses.ReadOnly...)

	for _, k := range res.accountKeys {
		res.accounts[k] = struct{}{}
	}

	res.tokenDeltas, err = NewBalanceDeltaMap(tx.Meta.PreTokenBalances, tx.Meta.PostTokenBalances, res.accountKeys)
	if err != nil {
		return nil, err
	}

	for i, inst := range res.parsed.Message.Instructions {
		res.flattenedInstructions = append(res.flattenedInstructions, inst)
		inner := res.GetInnerInstructions(uint16(i))
		res.flattenedInstructions = append(res.flattenedInstructions, inner...)
	}

	res.tokenTransfers, err = res.getTokenTransfers()
	if err != nil {
		return nil, err
	}

	res.solTransfers, err = res.getSolTransfers()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Reader) GetBlock() uint64 {
	return *r.block.BlockHeight
}

func (r *Reader) GetSlot() uint64 {
	return r.slot
}

func (r *Reader) GetTimestamp() time.Time {
	return r.block.BlockTime.Time()
}

func (r *Reader) GetTransactionIndex() int {
	return r.txIndex
}

func (r *Reader) GetSignature() solana.Signature {
	return r.parsed.Signatures[0]
}

func (r *Reader) CheckAccountKeysContains(addr solana.PublicKey) bool {
	_, ok := r.accounts[addr]
	return ok
}

func (r *Reader) GetAccountAtIndex(index uint16) (solana.PublicKey, error) {
	if int(index) > len(r.accountKeys) {
		return solana.PublicKey{}, errors.New("account index out of range")
	}
	return r.accountKeys[index], nil
}

func (r *Reader) GetAccountsAtIndices(indices []uint16) ([]solana.PublicKey, error) {
	accounts := make([]solana.PublicKey, 0, len(indices))
	for _, index := range indices {
		acc, err := r.GetAccountAtIndex(index)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acc)
	}

	return accounts, nil
}

func (r *Reader) AccountsLength() int {
	return len(r.accountKeys)
}

func (r *Reader) GetInstructions() []solana.CompiledInstruction {
	return r.parsed.Message.Instructions
}

func (r *Reader) GetInnerInstructions(topLevelInstructionIndex uint16) []solana.CompiledInstruction {
	res := make([]solana.CompiledInstruction, 0, len(r.tx.Meta.InnerInstructions))

	for _, inner := range r.tx.Meta.InnerInstructions {
		if inner.Index != uint16(topLevelInstructionIndex) {
			continue
		}

		res = append(res, inner.Instructions...)
	}

	return res
}

func (r *Reader) GetFlattenedInstructions() []solana.CompiledInstruction {
	return r.flattenedInstructions
}

func (r *Reader) GetInstructionAtFlattenedIndex(ind uint8) (solana.CompiledInstruction, error) {
	if ind >= uint8(len(r.flattenedInstructions)) {
		return solana.CompiledInstruction{}, fmt.Errorf("flattened instruction index out of range")
	}
	return r.flattenedInstructions[ind], nil
}

func (r *Reader) GetMintAndOwnerForTokenAccount(tokenAccount solana.PublicKey) (owner solana.PublicKey, mint solana.PublicKey, err error) {
	delta, ok := r.tokenDeltas.GetDelta(tokenAccount)
	if !ok {
		return owner, mint, errors.New("error retrieving token account")
	}

	return delta.Owner, delta.Mint, nil
}

func (r *Reader) GetTokenAccountBalanceChange(tokenAccount solana.PublicKey) (int64, error) {
	delta, ok := r.tokenDeltas.GetDelta(tokenAccount)
	if !ok {
		return 0, errors.New("error retrieving token account")
	}

	balanceAfter, err := strconv.ParseUint(delta.BalanceAfter.Amount, 10, 64)
	if err != nil {
		return 0, err
	}
	balanceBefore, err := strconv.ParseUint(delta.BalanceBefore.Amount, 10, 64)
	if err != nil {
		return 0, err
	}

	return int64(balanceAfter) - int64(balanceBefore), nil
}

func (r *Reader) GetInstructionAccounts(inst solana.CompiledInstruction) ([]*solana.AccountMeta, error) {
	return inst.ResolveInstructionAccounts(&r.parsed.Message)
}

func (r *Reader) GetParsedLogs() ([]logparse.Log, error) {
	return logparse.ParseLogs(r.tx.Meta.LogMessages)
}

type TokenTransfer struct {
	FlattenedInstructionIndex uint8
	Amount                    uint64
	Authority                 solana.PublicKey
	Destination               solana.PublicKey
	Source                    solana.PublicKey
}

func (r *Reader) GetTokenTransfers() []*TokenTransfer {
	return r.tokenTransfers
}

func (r *Reader) FilterTransfer(fn func(t *TokenTransfer) bool) ([]*TokenTransfer, error) {
	res := make([]*TokenTransfer, 0, len(r.tokenTransfers))
	for _, transfer := range r.tokenTransfers {
		if fn(transfer) {
			res = append(res, transfer)
		}
	}
	return res, nil
}

func (r *Reader) FindTransfer(fn func(t *TokenTransfer) bool) (*TokenTransfer, error) {
	for _, transfer := range r.tokenTransfers {
		if fn(transfer) {
			return transfer, nil
		}
	}

	return nil, fmt.Errorf("transfer not found")
}

func (r *Reader) FindInstruction(initialFlatInstIndex uint8, fn func(inst solana.CompiledInstruction, flatIndex uint8) bool) (solana.CompiledInstruction, error) {
	for i := int(initialFlatInstIndex); i < len(r.flattenedInstructions); i += 1 {
		next := r.flattenedInstructions[i]
		if fn(next, uint8(i)) {
			return next, nil
		}
	}
	return solana.CompiledInstruction{}, fmt.Errorf("instruction not found")
}

func (r *Reader) getTokenTransfers() ([]*TokenTransfer, error) {
	var res []*TokenTransfer
	for i, inst := range r.flattenedInstructions {
		programId, err := r.GetAccountAtIndex(inst.ProgramIDIndex)
		if err != nil {
			return nil, err
		}
		if !programId.Equals(solana.TokenProgramID) {
			continue
		}

		instType := uint8(0)
		err = bin.NewBinDecoder(inst.Data).Decode(&instType)
		if err != nil {
			return nil, err
		}

		if instType != 3 && instType != 12 {
			// Skip since the instruction is not a transfer
			continue
		}

		accounts, err := r.GetAccountsAtIndices(inst.Accounts)
		if err != nil {
			return nil, err
		}

		accountMeta := make([]*solana.AccountMeta, 0, len(accounts))
		for _, acc := range accounts {
			accountMeta = append(accountMeta, &solana.AccountMeta{PublicKey: acc})
		}

		parsed, err := token.DecodeInstruction(accountMeta, inst.Data)
		if err != nil {
			return nil, err
		}

		switch t := parsed.Impl.(type) {
		case *token.Transfer:
			res = append(res, &TokenTransfer{
				FlattenedInstructionIndex: uint8(i),
				Amount:                    *t.Amount,
				Authority:                 t.GetOwnerAccount().PublicKey,
				Source:                    t.GetSourceAccount().PublicKey,
				Destination:               t.GetDestinationAccount().PublicKey,
			})
		case *token.TransferChecked:
			res = append(res, &TokenTransfer{
				FlattenedInstructionIndex: uint8(i),
				Amount:                    *t.Amount,
				Authority:                 t.GetOwnerAccount().PublicKey,
				Source:                    t.GetSourceAccount().PublicKey,
				Destination:               t.GetDestinationAccount().PublicKey,
			})
		}
	}

	return res, nil
}

type SolTransfer struct {
	FlattenedInstructionIndex uint8
	Amount                    uint64
	Source                    solana.PublicKey
	Destination               solana.PublicKey
}

func (r *Reader) GetSolTransfers() []*SolTransfer {
	return r.solTransfers
}

func (r *Reader) getSolTransfers() ([]*SolTransfer, error) {
	var res []*SolTransfer
	for i, inst := range r.flattenedInstructions {
		programId, err := r.GetAccountAtIndex(inst.ProgramIDIndex)
		if err != nil {
			return nil, err
		}

		if !programId.Equals(solana.SystemProgramID) {
			continue
		}

		accounts, err := r.GetAccountsAtIndices(inst.Accounts)
		if err != nil {
			return nil, err
		}

		accountMeta := make([]*solana.AccountMeta, 0, len(accounts))
		for _, acc := range accounts {
			accountMeta = append(accountMeta, &solana.AccountMeta{PublicKey: acc})
		}

		parsed, err := system.DecodeInstruction(accountMeta, inst.Data)
		if err != nil {
			return nil, err
		}

		switch t := parsed.Impl.(type) {
		case *system.Transfer:
			res = append(res, &SolTransfer{
				FlattenedInstructionIndex: uint8(i),
				Amount:                    *t.Lamports,
				Source:                    t.GetFundingAccount().PublicKey,
				Destination:               t.GetRecipientAccount().PublicKey,
			})
		case *system.TransferWithSeed:
			res = append(res, &SolTransfer{
				FlattenedInstructionIndex: uint8(i),
				Amount:                    *t.Lamports,
				Source:                    t.GetFundingAccount().PublicKey,
				Destination:               t.GetRecipientAccount().PublicKey,
			})
		}
	}

	return res, nil
}

func (r *Reader) FilterSolTransfers(fn func(s *SolTransfer) bool) ([]*SolTransfer, error) {
	res := make([]*SolTransfer, 0, len(r.solTransfers))
	for _, transfer := range r.solTransfers {
		if fn(transfer) {
			res = append(res, transfer)
		}
	}
	return res, nil
}

func (r *Reader) FindSolTransfer(fn func(t *SolTransfer) bool) (*SolTransfer, error) {
	for _, transfer := range r.solTransfers {
		if fn(transfer) {
			return transfer, nil
		}
	}

	return nil, fmt.Errorf("transfer not found")
}

func (r *Reader) GetSolBalanceDelta(accountIndex uint16) (int64, error) {
	if int(accountIndex) > len(r.tx.Meta.PreBalances) || int(accountIndex) > len(r.tx.Meta.PostBalances) {
		return 0, fmt.Errorf("account index out of bounds")
	}
	preBalance := int64(r.tx.Meta.PreBalances[accountIndex])
	postBalance := int64(r.tx.Meta.PostBalances[accountIndex])
	return postBalance - preBalance, nil
}
