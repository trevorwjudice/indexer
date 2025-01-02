package keycache

import (
	"context"

	"github.com/gagliardetto/solana-go"
	"github.com/upper/db/v4"
)

type FetchFunc func(ctx context.Context, s db.Session) ([]solana.PublicKey, error)

type Keycache struct {
	fetch FetchFunc
	keys  map[solana.PublicKey]struct{}
}

func New(fn FetchFunc) *Keycache {
	return &Keycache{
		fetch: fn,
		keys:  make(map[solana.PublicKey]struct{}),
	}
}

func (k *Keycache) Fetch(ctx context.Context, s db.Session) error {
	keys, err := k.fetch(ctx, s)
	if err != nil {
		return err
	}
	for _, key := range keys {
		k.keys[key] = struct{}{}
	}
	return nil
}

func (k *Keycache) Add(key solana.PublicKey) bool {
	_, ok := k.keys[key]
	if ok {
		return true
	}
	k.keys[key] = struct{}{}
	return false
}

func (k *Keycache) Remove(key solana.PublicKey) bool {
	_, ok := k.keys[key]
	if ok {
		delete(k.keys, key)
		return true
	}
	return false
}

func (k *Keycache) Contains(key solana.PublicKey) bool {
	_, ok := k.keys[key]
	return ok
}
