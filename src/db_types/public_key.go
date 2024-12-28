package db_types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

type SolanaPublicKey solana.PublicKey

func (p SolanaPublicKey) MarshalDB() (interface{}, error) {
	return solana.PublicKey(p).Bytes(), nil
}

func (p *SolanaPublicKey) UnmarshalDB(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot convert %T to []byte", value)
	}

	*p = SolanaPublicKey(solana.PublicKeyFromBytes(bytes))
	return nil
}

func (p SolanaPublicKey) MarshalJSON() ([]byte, error) {
	encoded := base58.Encode(solana.PublicKey(p).Bytes())
	return json.Marshal(encoded) // Return it as a JSON string
}

func (p *SolanaPublicKey) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return fmt.Errorf("failed to unmarshal public key from JSON: %w", err)
	}

	bytes, err := base58.Decode(encoded)
	if err != nil {
		return fmt.Errorf("failed to decode Base58 public key: %w", err)
	}

	*p = SolanaPublicKey(solana.PublicKeyFromBytes(bytes))
	return nil
}

func (p SolanaPublicKey) Value() (driver.Value, error) {
	return solana.PublicKey(p).Bytes(), nil // Return raw 32-byte binary data
}

func (p SolanaPublicKey) PublicKey() solana.PublicKey {
	return solana.PublicKey(p)
}
