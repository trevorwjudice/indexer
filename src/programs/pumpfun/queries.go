package pumpfun

import (
	"context"
	"indexer/src/db_types"
	"indexer/src/query"

	"github.com/upper/db/v4"
)

func GetCreateInstructions(ctx context.Context, s db.Session) ([]*db_types.PumpfunCreate, error) {
	return query.SqlQuery[*db_types.PumpfunCreate](ctx, s, func(_ context.Context, d db.SQL) db.ResultMapper {
		return d.SelectFrom(db_types.PUMP_FUN_CREATE)
	})
}
