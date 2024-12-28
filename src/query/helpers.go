package query

import (
	"context"

	"github.com/upper/db/v4"
)

func sqlQuery[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.ResultMapper) ([]T, error) {
	var o []T
	err := fn(ctx, s.SQL()).All(&o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func sqlQueryRow[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.ResultMapper) (T, error) {
	var val T
	err := fn(ctx, s.SQL()).One(&val)
	if err != nil {
		return val, err
	}
	return val, err
}

func sqlQueryType[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.Selector) ([]T, error) {
	var vals []T
	res, err := fn(ctx, s.SQL()).QueryContext(ctx)
	if err != nil {
		return vals, err
	}

	for res.Next() {
		v := new(T)
		err := res.Scan(v)
		if err != nil {
			return vals, err
		}
		vals = append(vals, *v)
	}

	return vals, nil
}

func sqlQueryTypeRow[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.Selector) (T, error) {
	val := new(T)
	res, err := fn(ctx, s.SQL()).QueryContext(ctx)
	if err != nil {
		return *val, err
	}
	res.Next()
	return *val, res.Scan(val)
}

func SqlQuery[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.ResultMapper) ([]T, error) {
	return sqlQuery[T](ctx, s, fn)
}

func SqlQueryRow[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.ResultMapper) (T, error) {
	return sqlQueryRow[T](ctx, s, fn)
}

func SqlQueryType[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.Selector) ([]T, error) {
	return sqlQueryType[T](ctx, s, fn)
}

func SqlQueryTypeRow[T any](ctx context.Context, s db.Session, fn func(context.Context, db.SQL) db.Selector) (T, error) {
	return sqlQueryTypeRow[T](ctx, s, fn)
}
