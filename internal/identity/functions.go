package identity

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"reflect"
	"sort"
	"strings"

	"github.com/jarrodhroberson/go-set/internal"
)

func includeFieldPredicate(f reflect.StructField, v reflect.Value) (bool, error) {
	if str := f.Tag.Get("identity"); str != "" {
		if str == "-" {
			return false, nil
		}
	}
	return true, nil
}

func primitiveStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	switch rv.Kind() {
	case reflect.String:
		return []byte(rv.String()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var b bytes.Buffer
		_ = binary.Write(&b, binary.BigEndian, rv.Int())
		return b.Bytes(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var b bytes.Buffer
		_ = binary.Write(&b, binary.BigEndian, rv.Uint())
		return b.Bytes(), true
	case reflect.Float32, reflect.Float64:
		var b bytes.Buffer
		_ = binary.Write(&b, binary.BigEndian, rv.Float())
		return b.Bytes(), true
	case reflect.Bool:
		var b bytes.Buffer
		_ = binary.Write(&b, binary.BigEndian, rv.Bool())
		return b.Bytes(), true
	default:
		return nil, false
	}
}

func pointerStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	if rv.Kind() == reflect.Ptr {
		if !rv.IsNil() || rv.Type().Elem().Kind() == reflect.Struct {
			return structStrategy(h, rv)
		} else {
			zero := reflect.Zero(rv.Type().Elem())
			if b, ok := primitiveStrategy(h, zero); ok {
				return b, true
			} else if b, ok := mapStrategy(h, zero); ok {
				return b, true
			} else if b, ok := pointerStrategy(h, zero); ok {
				return b, true
			}
			return nil, false
		}
	}
	return nil, false
}

func mapStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	if rv.Kind() == reflect.Map {
		mk := rv.MapKeys()
		kv := make(map[string]reflect.Value, len(mk))
		for _, k := range mk {
			kv[k.String()] = k
		}
		keys := internal.MapKeysAsSlice[string, reflect.Value](kv)
		sort.Strings(keys)
		b := bytes.Buffer{}
		for idx := range keys {
			strategies{
				primitiveStrategy,
				pointerStrategy,
				mapStrategy,
				structStrategy,
				interfaceStrategy,
				arraySliceStrategy,
				defaultStrategy,
			}.apply(h, rv.MapIndex(kv[keys[idx]]))
		}
		return b.Bytes(), true
	}
	return nil, false
}

func structStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	if rv.Kind() == reflect.Struct {
		vtype := rv.Type()
		flen := vtype.NumField()
		kv := make(map[string]reflect.Value, flen)

		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			ok, err := includeFieldPredicate(field, rv.Field(i))
			if err != nil && strings.Contains(err.Error(), "method:") {
				panic(err)
			}
			if !ok {
				continue
			}
			kv[field.Name] = rv.Field(i)
		}

		keys := internal.MapKeysAsSlice[string, reflect.Value](kv)
		sort.Strings(keys)
		b := bytes.Buffer{}
		for idx := range keys {
			strategies{
				primitiveStrategy,
				pointerStrategy,
				mapStrategy,
				structStrategy,
				interfaceStrategy,
				arraySliceStrategy,
				defaultStrategy,
			}.apply(h, rv.MapIndex(kv[keys[idx]]))
		}
		return b.Bytes(), true
	}
	return nil, false
}

func interfaceStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	if rv.Kind() == reflect.Interface {
		if !rv.CanInterface() {
			return nil, false
		}
		strategies{
			primitiveStrategy,
			pointerStrategy,
			mapStrategy,
			structStrategy,
			interfaceStrategy,
			arraySliceStrategy,
			defaultStrategy,
		}.apply(h, reflect.ValueOf(rv.Interface()))
	}
	return nil, false
}

func arraySliceStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		var b bytes.Buffer
		for i := 0; i < rv.Len(); i++ {
			strategies{
				primitiveStrategy,
				pointerStrategy,
				mapStrategy,
				structStrategy,
				interfaceStrategy,
				arraySliceStrategy,
				defaultStrategy,
			}.apply(h, rv)
		}
		return b.Bytes(), true
	default:
		return nil, false
	}
}

func defaultStrategy(h hash.Hash, rv reflect.Value) ([]byte, bool) {
	return rv.Bytes(), true
}

type strategies []func(h hash.Hash, rv reflect.Value) ([]byte, bool)

var identityStrategies = strategies{
	primitiveStrategy,
	pointerStrategy,
	mapStrategy,
	structStrategy,
	interfaceStrategy,
	arraySliceStrategy,
	defaultStrategy,
}

func (is strategies) apply(h hash.Hash, object any) []byte {
	for _, strategy := range is {
		if b, ok := strategy(h, reflect.ValueOf(object)); ok {
			return b
		}
	}
	return []byte{}
}

func hashAny[T any](object T, h hash.Hash) []byte {
	return h.Sum(identityStrategies.apply(h, object))
}

func HashIdentity[T any](t T) string {
	return hex.EncodeToString(hashAny(t, sha512.New()))
}
