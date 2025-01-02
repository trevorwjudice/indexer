package queries

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/query"

	"github.com/upper/db/v4"
)

func GetPumpFunCreateInstructions(ctx context.Context, s db.Session) ([]*db_types.PumpFunCreate, error) {
	return query.SqlQuery[*db_types.PumpFunCreate](ctx, s, func(_ context.Context, d db.SQL) db.ResultMapper {
		return d.SelectFrom(db_types.PUMP_FUN_CREATE)
	})
}

func GetPumpFunSwaps(ctx context.Context, s db.Session) ([]*db_types.PumpFunSwap, error) {
	return query.SqlQuery[*db_types.PumpFunSwap](ctx, s, func(_ context.Context, d db.SQL) db.ResultMapper {
		return d.SelectFrom(db_types.PUMP_FUN_SWAPS)
	})
}
