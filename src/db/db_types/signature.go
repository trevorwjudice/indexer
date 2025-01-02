package db_types

import (
	"database/sql/driver"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

type Signature solana.Signature

func (p Signature) MarshalDB() (interface{}, error) {
	b := [64]byte(p.Signature())
	return b[:], nil
}

func (p *Signature) UnmarshalDB(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %T to []byte", value)
	}
	*p = Signature(solana.SignatureFromBytes(bytes))
	return nil
}

func (p Signature) MarshalJSON() ([]byte, error) {
	return solana.Signature(p).MarshalJSON()
}

func (p *Signature) UnmarshalJSON(data []byte) error {
	var sig solana.Signature

	err := sig.UnmarshalJSON(data)
	if err != nil {
		return err
	}

	*p = Signature(sig)
	return nil
}

func (p Signature) Value() (driver.Value, error) {
	return solana.Signature(p).MarshalJSON() // Return raw 32-byte binary data
}

func (p Signature) Signature() solana.Signature {
	return solana.Signature(p)
}
