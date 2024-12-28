package dbutil

import "github.com/upper/db/v4"

func BatchUpload[T any](s db.Session, tbl string, items []T) error {
	inserter := s.SQL().InsertInto(tbl).Amend(func(q string) string {
		return q + " ON CONFLICT DO NOTHING "
	}).Batch(1000)
	go func() {
		defer inserter.Done()
		for _, item := range items {
			inserter.Values(item)
		}
	}()
	return inserter.Wait()
}
