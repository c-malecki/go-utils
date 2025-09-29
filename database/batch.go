package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/c-malecki/go-utils/parse/pslice"
)

type BatchInsertDesc[T any] struct {
	Query     string
	Items     []T
	ExtractFn func(T) []interface{}
}

type BatchInsertResult[T any] struct {
	ID     int64
	Entity T
}

func BatchInsert[T any](ctx context.Context, db *sql.DB, desc BatchInsertDesc[T]) ([]BatchInsertResult[T], error) {
	if len(desc.Items) == 0 {
		return []BatchInsertResult[T]{}, nil
	}
	results := make([]BatchInsertResult[T], 0, len(desc.Items))
	parts := strings.SplitAfter(desc.Query, "VALUES")
	if len(parts) != 2 {
		return nil, fmt.Errorf("missing VALUES in insert query")
	}
	base := parts[0]
	bindvars := parts[1]
	count := strings.Count(bindvars, "?")
	// max bind variables is 65,536 but using 40000 to cut down on batch sizing
	size := 40000 / count

	split := pslice.SubSlice(desc.Items, size)
	for _, sub := range split {
		tx, err := db.Begin()
		if err != nil {
			return nil, fmt.Errorf("db.Begin: %w", err)
		}
		defer tx.Rollback()

		placeholders := make([]string, 0)
		args := make([]interface{}, 0)

		for _, v := range sub {
			placeholders = append(placeholders, bindvars)
			a := desc.ExtractFn(v)
			args = append(args, a...)
		}

		query := base + strings.Join(placeholders, ",")

		res, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("tx.ExecContext: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("tx.Commit: %w", err)
		}

		firstId, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("res.LastInsertId: %w", err)
		}

		for i, v := range sub {
			id := int64(firstId) + int64(i)
			results = append(results, BatchInsertResult[T]{
				ID:     id,
				Entity: v,
			})
		}
	}

	return results, nil
}

func BatchInsertWithTx[T any](ctx context.Context, tx *sql.Tx, desc BatchInsertDesc[T]) ([]BatchInsertResult[T], error) {
	if len(desc.Items) == 0 {
		return []BatchInsertResult[T]{}, nil
	}
	results := make([]BatchInsertResult[T], 0, len(desc.Items))
	parts := strings.SplitAfter(desc.Query, "VALUES")
	if len(parts) != 2 {
		return nil, fmt.Errorf("missing VALUES in insert query")
	}
	base := parts[0]
	bindvars := parts[1]
	count := strings.Count(bindvars, "?")
	// max bind variables is 65,536 but using 40000 to cut down on batch sizing
	size := 40000 / count

	split := pslice.SubSlice(desc.Items, size)
	for _, sub := range split {
		placeholders := make([]string, 0)
		args := make([]interface{}, 0)

		for _, v := range sub {
			placeholders = append(placeholders, bindvars)
			a := desc.ExtractFn(v)
			args = append(args, a...)
		}

		query := base + strings.Join(placeholders, ",")

		res, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("tx.ExecContext: %w", err)
		}

		firstId, err := res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("res.LastInsertId: %w", err)
		}

		for i, v := range sub {
			id := int64(firstId) + int64(i)
			results = append(results, BatchInsertResult[T]{
				ID:     id,
				Entity: v,
			})
		}
	}

	return results, nil
}

func InsertManyAndReturnIDsWithTx(ctx context.Context, tx *sql.Tx, query string, args []interface{}) ([]int64, error) {
	createdIds := make([]int64, 0)

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return createdIds, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return createdIds, err
	}

	firstId, err := res.LastInsertId()
	if err != nil {
		return createdIds, err
	}

	lastId := firstId + (affected - 1)
	curId := firstId

	for curId <= lastId {
		id := int64(curId)
		createdIds = append(createdIds, id)
		curId += 1
	}

	return createdIds, nil
}

func InsertManyAndReturnIDs(ctx context.Context, db *sql.DB, query string, args []interface{}) ([]int64, error) {
	createdIds := make([]int64, 0)

	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return createdIds, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return createdIds, err
	}

	firstId, err := res.LastInsertId()
	if err != nil {
		return createdIds, err
	}

	lastId := firstId + (affected - 1)
	curId := firstId

	for curId <= lastId {
		id := int64(curId)
		createdIds = append(createdIds, id)
		curId += 1
	}

	return createdIds, nil
}
