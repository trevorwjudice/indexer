package indexer

import (
	"context"
	"indexer/src/db/db_types"
	"indexer/src/query"

	"github.com/upper/db/v4"
)

func (i *Indexer) resetTaskProgress(ctx context.Context) error {
	_, err := i.s.SQL().Update(db_types.INDEXER_PROGRESS).Set("status", 0).Where("status = -1").ExecContext(ctx)
	return err
}

func (i *Indexer) getScheduledSlotHeight(ctx context.Context) (uint64, error) {
	type scheduledSlotHeight struct {
		Slot uint64 `db:"slot_end"`
	}
	height := &scheduledSlotHeight{}
	err := i.s.SQL().Select("slot_end").From(db_types.INDEXER_PROGRESS).OrderBy("-slot_end").Limit(1).IteratorContext(ctx).One(height)
	if err != nil {
		return 0, nil
	}
	return height.Slot, nil
}

func (i *Indexer) startNextTask(ctx context.Context) (v *db_types.Progress, err error) {
	return v, i.s.TxContext(ctx, func(d db.Session) error {
		v, err = query.SqlQueryRow[*db_types.Progress](ctx, i.s, func(_ context.Context, d db.SQL) db.ResultMapper {
			return d.SelectFrom(db_types.INDEXER_PROGRESS).Where("status = 0").OrderBy("slot_start").Limit(1)
		})
		if err != nil {
			return err
		}
		v.Status = -1
		return d.Collection(db_types.INDEXER_PROGRESS).UpdateReturning(v)
	}, nil)
}

func (i *Indexer) finishTask(s db.Session, t *db_types.Progress) error {
	return s.Collection(db_types.INDEXER_PROGRESS).UpdateReturning(t)
}
