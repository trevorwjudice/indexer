package queries

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/query"

	"github.com/upper/db/v4"
)

func GetRaydiumSwaps(ctx context.Context, s db.Session, limit int) ([]*db_types.RaydiumSwap, error) {
	return query.SqlQuery[*db_types.RaydiumSwap](ctx, s, func(_ context.Context, d db.SQL) db.ResultMapper {
		return d.SelectFrom(db_types.RAYDIUM_V4_SWAPS).Limit(limit)
	})
}
