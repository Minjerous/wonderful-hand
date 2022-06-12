package helper

import (
	"context"
	"database/sql"
	"github.com/igxnon/cachepool"
	"github.com/igxnon/cachepool/helper/internal"
)

type ExecResult struct {
	Result sql.Result
	Err    error
}

// Query try to get data in cache and return generic type slice []T, if data is not found in cache
// then go to SQL database to get
// the param key set into cache and used to locate cache in the next Query
// generic type T support map[string]any, []map[string]any, struct, []struct,
// int, []int, string, []string, float, []float, bool, []bool
func Query[T any](c cachepool.ICachePool, key, query string, args ...any) (rows []T, err error) {
	return QueryWithContext[T](context.Background(), c, key, query, args...)
}

// QueryWithContext try to get data in cache and return generic type slice []T, if data is not found in cache
// then go to SQL database to get
// the param key set into cache and used to locate cache in the next QueryWithContext
// generic type T support map[string]any, []map[string]any, struct, []struct,
// int, []int, string, []string, float, []float, bool, []bool
func QueryWithContext[T any](
	ctx context.Context,
	c cachepool.ICachePool,
	key, query string, args ...any,
) (rows []T, err error) {
	return internal.HandleRows[[]T](ctx, c, key, query, args...)
}

// QueryRow try to get data in cache and return generic type T, if data is not found in cache
// then go to SQL database to get
// the param key set into cache and used to locate cache in the next QueryRow
// generic type T support map[string]any, []map[string]any, struct, []struct,
// int, []int, string, []string, float, []float, bool, []bool
func QueryRow[T any](c cachepool.ICachePool, key, query string, args ...any) (rows T, err error) {
	return QueryRowWithContext[T](context.Background(), c, key, query, args...)
}

// QueryRowWithContext try to get data in cache and return generic type T, if data is not found in cache
// then go to SQL database to get
// the param key set into cache and used to locate cache in the next QueryRowWithContext
// generic type T support map[string]any, []map[string]any, struct, []struct,
// int, []int, string, []string, float, []float, bool, []bool
func QueryRowWithContext[T any](
	ctx context.Context,
	c cachepool.ICachePool,
	key, query string, args ...any,
) (row T, err error) {
	return internal.HandleRow[T](ctx, c, key, query, args...)
}
