package db_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

type PublicKey solana.PublicKey

func (p PublicKey) MarshalDB() (interface{}, error) {
	return solana.PublicKey(p).Bytes(), nil
}

func (p *PublicKey) UnmarshalDB(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %T to []byte", value)
	}
	*p = PublicKey(solana.PublicKeyFromBytes(bytes))
	return nil
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	encoded := base58.Encode(solana.PublicKey(p).Bytes())
	return json.Marshal(encoded) // Return it as a JSON string
}

func (p *PublicKey) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return fmt.Errorf("failed to unmarshal public key from JSON: %w", err)
	}

	bytes, err := base58.Decode(encoded)
	if err != nil {
		return fmt.Errorf("failed to decode Base58 public key: %w", err)
	}

	*p = PublicKey(solana.PublicKeyFromBytes(bytes))
	return nil
}

func (p PublicKey) Value() (driver.Value, error) {
	return solana.PublicKey(p).Bytes(), nil // Return raw 32-byte binary data
}

func (p PublicKey) PublicKey() solana.PublicKey {
	return solana.PublicKey(p)
}

func ToPublicKeyErr(k solana.PublicKey, err error) (PublicKey, error) {
	res := PublicKey(k)
	return res, err
}
