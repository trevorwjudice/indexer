package db_types

import (
	"database/sql/driver"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type SignatureWrapper solana.Signature

func (p SignatureWrapper) MarshalDB() (interface{}, error) {
	return solana.Signature(p).MarshalJSON()
}

func (p *SignatureWrapper) UnmarshalDB(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %T to []byte", value)
	}

	*p = SignatureWrapper(solana.SignatureFromBytes(bytes))
	return nil
}

func (p SignatureWrapper) MarshalJSON() ([]byte, error) {
	return solana.Signature(p).MarshalJSON()
}

func (p *SignatureWrapper) UnmarshalJSON(data []byte) error {
	var sig solana.Signature

	err := sig.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*p = SignatureWrapper(sig)
	return nil
}

func (p SignatureWrapper) Value() (driver.Value, error) {
	return solana.Signature(p).MarshalJSON() // Return raw 32-byte binary data
}

func (p SignatureWrapper) Signature() solana.Signature {
	return solana.Signature(p)
}
