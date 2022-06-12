package helper

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/igxnon/cachepool"
	"testing"
)

type FooBar struct {
	Bar int64
	Yee string
	Foo sql.NullTime
}

type BadFooBar struct {
	Foo int64
	Yee string
	Bar sql.NullTime
}

/*
------------------
yee	  | bar	| foo
------------------
Hello |  1  | null
------------------
Hi	  |  2  | NOW()
------------------
*/

func TestQuery(t *testing.T) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Error(err)
		return
	}
	pool := cachepool.New(cachepool.WithDatabase(db))
	got, err := QueryRow[FooBar](pool, "foobar:combine", "SELECT * FROM t LIMIT 1 OFFSET 1")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v\n", got)

	gots, err := Query[map[string]any](pool, "foobar:combine", "SELECT * FROM t LIMIT 5")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(gots)

	gotOnes, err := Query[int32](pool, "bar:int", "SELECT bar FROM t LIMIT 5")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(gotOnes)

	gotOnesNullable, err := Query[sql.NullTime](pool, "foo:time", "SELECT foo FROM t LIMIT 5")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(gotOnesNullable)
}

func TestBadQuery(t *testing.T) {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Error(err)
		return
	}
	pool := cachepool.New(cachepool.WithDatabase(db))
	got, err := QueryRow[BadFooBar](pool, "foobar:combine", "SELECT * FROM t WHERE bar = ? LIMIT 1", 1)
	if err == nil {
		t.Error("opps")
		return
	}
	t.Logf("%#v, %v\n", got, err)

	gotOnes, err := Query[sql.NullTime](pool, "bar:int", "SELECT bar FROM t LIMIT 5")
	if err == nil {
		t.Error("opps")
		return
	}
	t.Logf("%#v, %v\n", gotOnes, err)
}

func BenchmarkQueryWithCache(b *testing.B) {
	var (
		out     map[string]any
		dsn     = "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4"
		db, err = sql.Open("mysql", dsn)
		pool    = cachepool.New(cachepool.WithDatabase(db))
	)

	if err != nil {
		b.Error(err)
	}

	_, err = QueryRow[map[string]any](pool, "foobar:combine", "SELECT * FROM t WHERE bar = ? LIMIT 1", 1)
	if err != nil {
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err = QueryRow[map[string]any](pool, "foobar:combine", "SELECT * FROM t WHERE bar = ? LIMIT 1", 1)
		if err != nil {
			b.Error(err)
		}
		if i%500000000 == 0 {
			b.Log(out)
		}
	}
}

func BenchmarkQueryWithoutCache(b *testing.B) {
	var (
		dsn     = "root:12345678@tcp(127.0.0.1:3306)/awesome?charset=utf8mb4"
		db, err = sql.Open("mysql", dsn)
	)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var yee string
		var bar int
		var foo sql.NullTime
		row := db.QueryRow("SELECT * FROM t WHERE bar = ? LIMIT 1", 1)
		err = row.Scan(&yee, &bar, &foo)
		if err != nil {
			b.Error(err)
		}
		if i%50000 == 0 {
			b.Log(foo, bar)
		}
	}
}
