package internal

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/alecthomas/binary"
	"github.com/igxnon/cachepool"
	"reflect"
	"time"
	"unicode"
)

type column struct {
	typ       reflect.Type
	fieldName string
	fieldPos  int
}

// type map[string]any, []map[string]any, struct, []struct,
// int, []int, string, []string, float, []float, bool, []bool are supported
func check(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Map:
		return typ.Key().Kind() == reflect.String
	case reflect.Struct:
		return true
	case reflect.Slice:
		switch elem := typ.Elem(); elem.Kind() {
		case reflect.Map:
			return elem.Key().Kind() == reflect.String
		case reflect.Struct:
			return true
		default:
			return checkElemType(typ.Elem()) == nil
		}
	default:
		return checkElemType(typ) == nil
	}
}

func parseStruct(typ reflect.Type) (map[string]column, error) {
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("parse failed, not a struct")
	}
	n := typ.NumField()
	m := make(map[string]column, n)
	for i := 0; i < n; i++ {
		var (
			field   = typ.Field(i)
			colname = camel2Case(field.Name)
		)
		if err := checkElemType(field.Type); err != nil {
			return nil, fmt.Errorf("struct field %s, %v", field.Name, err)
		}
		m[colname] = column{
			typ:       field.Type,
			fieldName: field.Name,
			fieldPos:  i,
		}
	}
	return m, nil
}

func camel2Case(name string) string {
	buffer := new(bytes.Buffer)
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func checkElemType(typ reflect.Type) error {
	if typ.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) ||
		reflect.PtrTo(typ).Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
		return nil
	}
	if typ.String() == reflect.TypeOf(time.Time{}).String() {
		return nil
	}
	switch typ.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64:
		return nil
	case reflect.Slice: // []byte
		if typ.Elem().Kind() == reflect.Uint8 {
			return nil
		}
		fallthrough
	default:
		return fmt.Errorf("unsupported field type %s\n", typ.String())
	}
}

func HandleRows[S ~[]E, E any](
	ctx context.Context,
	c cachepool.ICachePool,
	key, query string, args ...any,
) (s S, err error) {
	var (
		elem     = reflect.TypeOf((*E)(nil)).Elem()
		typ      = reflect.TypeOf((*S)(nil)).Elem()
		ok       bool
		r        *sql.Rows
		cols     []string
		coltypes []*sql.ColumnType
	)

	if !check(typ) {
		err = errors.New("unsupported generic type")
		return
	}

	if err != nil {
		err = errors.New("sql syntax fault")
		return
	}

	s, ok = checkCache[S](key, c)
	if ok {
		return
	}

	// missed, go to database to get

	r, cols, coltypes, err = queryDb(c.GetDatabase(), ctx, query, args...)
	if err != nil {
		return
	}

	for r.Next() {
		var (
			values = make([]any, len(coltypes))
			e      E
		)
		err = scan(r, coltypes, values)
		if err != nil {
			return
		}
		e, err = bind[E](elem, values, cols)
		s = append(s, e)
	}

	saveToCache(c, key, s)
	return
}

func HandleRow[E any](
	ctx context.Context,
	c cachepool.ICachePool,
	key, query string, args ...any,
) (e E, err error) {
	var (
		typ      = reflect.TypeOf((*E)(nil)).Elem()
		r        *sql.Rows
		cols     []string
		coltypes []*sql.ColumnType
		ok       bool
	)

	if err != nil {
		err = errors.New("sql syntax fault")
		return
	}

	if !check(typ) {
		err = errors.New("unsupported generic type")
		return
	}

	e, ok = checkCache[E](key, c)
	if ok {
		return
	}

	// missed, go to database to get
	r, cols, coltypes, err = queryDb(c.GetDatabase(), ctx, query, args...)
	if err != nil {
		return
	}

	if !r.Next() {
		err = fmt.Errorf("row has no next, err: %v", r.Err())
		return
	}

	values := make([]any, len(coltypes))
	err = scan(r, coltypes, values)
	if err != nil {
		return
	}

	e, err = bind[E](typ, values, cols)
	if err != nil {
		return
	}

	saveToCache(c, key, e)
	return
}

func prepareValue(values []any, coltypes []*sql.ColumnType) {
	for idx, columnType := range coltypes {
		if columnType.ScanType() != nil {
			values[idx] = reflect.New(reflect.PtrTo(columnType.ScanType())).Interface()
		} else {
			values[idx] = new(interface{})
		}
	}
}

func writeToMap(value []any, cols []string, mapValue map[string]any) {
	for idx, col := range cols {
		if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(value[idx]))); reflectValue.IsValid() {
			mapValue[col] = reflectValue.Interface()
			if valuer, ok := mapValue[col].(driver.Valuer); ok {
				mapValue[col], _ = valuer.Value()
			} else if b, ok := mapValue[col].(sql.RawBytes); ok {
				mapValue[col] = string(b)
			}
		} else {
			mapValue[col] = nil
		}
	}
}

func writeToStruct[E any](value []any, cols []string, columns map[string]column, e *E) (err error) {
	for idx, col := range cols {
		var (
			c, ok = columns[col]
			field = reflect.ValueOf(e).Elem().Field(c.fieldPos)
			dest  = field.Interface()
		)
		if !ok {
			continue
		}

		dest, err = convert(value[idx], field.Type())
		if err != nil {
			return
		}

		if dest != nil {
			field.Set(reflect.ValueOf(dest))
		}
	}
	return
}

func saveToCache(c cachepool.ICachePool, key string, data any) {
	// TODO customized expire time
	c.SetDefault(key, data)
}

// for sugar global cache
type unmarshalable interface {
	GetUnmarshal(k string, obj interface{}) bool
}

func checkCache[T any](key string, c cachepool.ICachePool) (T, bool) {
	var t T
	got, ok := c.Get(key)
	if ok {
		t, ok = got.(T)
		if ok {
			return t, true
		}
		if _, ok = c.GetImplementedCache().(unmarshalable); ok {
			if g, ok := got.([]byte); ok {
				ok = binary.Unmarshal(g, &t) == nil
				return t, ok
			}
		}
		c.Delete(key)
		return t, false
	}
	return t, false
}

func queryDb(db *sql.DB, ctx context.Context, query string, args ...any) (r *sql.Rows, cols []string, coltypes []*sql.ColumnType, err error) {
	if db == nil {
		err = errors.New("database not found")
		return
	}
	r, err = db.QueryContext(ctx, query, args...)
	if err != nil {
		err = errors.New("query to database failed")
		return
	}
	cols, err = r.Columns()
	if err != nil {
		err = errors.New("get query columns failed")
		return
	}
	coltypes, err = r.ColumnTypes()
	if err != nil {
		err = errors.New("get query columns types failed")
		return
	}
	return
}

func scan(r *sql.Rows, coltypes []*sql.ColumnType, values []any) (err error) {
	prepareValue(values, coltypes)
	err = r.Scan(values...)
	if err != nil {
		err = errors.New("scan rows failed")
		return
	}
	return
}

type CastType uint

const (
	Map CastType = iota
	Struct
	Single
)

func bind[E any](typ reflect.Type, values []any, cols []string) (e E, err error) {
	switch figureType(typ) {
	case Map:
		var (
			mapValue any = make(map[string]any, len(cols))
			ok       bool
		)
		writeToMap(values, cols, mapValue.(map[string]any))
		e, ok = mapValue.(E)
		if !ok {
			err = errors.New("return value cannot cast to E, if used map as generic type make sure it equals map[string]any")
			return
		}
	case Struct:
		var columns map[string]column
		columns, err = parseStruct(typ)
		if err != nil {
			return
		}
		err = writeToStruct(values, cols, columns, &e)
		if err != nil {
			return
		}
	case Single:
		if len(values) != 1 {
			err = errors.New("selected value length not equal 1, please check your sql syntax")
			return
		}
		var (
			v    = values[0]
			dest any
			ok   bool
		)
		dest, err = convert(v, typ)
		if dest != nil {
			e, ok = dest.(E)
			if !ok {
				err = errors.New("return value cannot cast to E, if used map as generic type make sure it equals map[string]any")
				return
			}
		}
		return
	}
	return
}

func figureType(typ reflect.Type) CastType {
	switch typ.Kind() {
	case reflect.Map:
		return Map
	case reflect.Struct:
		if typ.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) ||
			reflect.PtrTo(typ).Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) ||
			typ.String() == reflect.TypeOf(time.Time{}).String() {
			return Single
		}
		return Struct
	default:
		return Single
	}
}

func convert(sv any, typ reflect.Type) (e any, err error) {
	if reflectValue := reflect.Indirect(reflect.Indirect(reflect.ValueOf(sv))); reflectValue.IsValid() {
		v := reflectValue.Interface()
		if valuer, ok := v.(driver.Valuer); ok {
			v, _ = valuer.Value()
		} else if b, ok := v.(sql.RawBytes); ok {
			v = string(b)
		}

		if reflect.TypeOf(v).ConvertibleTo(typ) {
			// v type can cast to field type
			e = reflect.ValueOf(v).Convert(typ).Interface()
			return
		} else if typ.Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
			// dest implement sql.Scanner
			err = e.(sql.Scanner).Scan(v)
			return
		} else if reflect.PtrTo(typ).Implements(reflect.TypeOf((*sql.Scanner)(nil)).Elem()) {
			// *dest implement sql.Scanner
			scanner := reflect.New(typ).Interface().(sql.Scanner)
			err = scanner.Scan(v)
			if err != nil {
				return
			}
			reflect.ValueOf(&e).Elem().Set(reflect.ValueOf(scanner).Elem())
			return
		}
		err = fmt.Errorf("convert value to type %s failed", typ.String())
		return
	}
	// sv equals nil skip it
	return e, nil
}
