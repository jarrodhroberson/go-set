package identity

import (
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"reflect"
	"sort"
	"strings"
)

func keysToSlice[K comparable, V any](m map[K]V) []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}

func includeFieldPredicate(f reflect.StructField, v reflect.Value) (bool, error) {
	if str := f.Tag.Get("identity"); str != "" {
		if str == "-" {
			return false, nil
		}
	}
	return true, nil
}

func primitiveStrategy(h hash.Hash, rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.String:
		h.Sum([]byte(rv.String()))
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		h.Sum()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func pointerStrategy(h hash.Hash, rv reflect.Value) bool {
	if rv.Kind() == reflect.Ptr {
		if !rv.IsNil() || rv.Type().Elem().Kind() == reflect.Struct {
			return structStrategy(h, rv)
		} else {
			zero := reflect.Zero(rv.Type().Elem())
			if primitiveStrategy(h, zero) {
				return true
			} else if mapStrategy(h, zero) {
				return true
			} else if pointerStrategy(h, zero) {
				return true
			}
			return false
		}
	}
	return false
}

func mapStrategy(h hash.Hash, rv reflect.Value) bool {
	if rv.Kind() == reflect.Map {
		mk := rv.MapKeys()
		kv := make(map[string]reflect.Value, len(mk))
		for _, k := range mk {
			kv[k.String()] = k
		}
		keys := keysToSlice[string, reflect.Value](kv)
		sort.Strings(keys)
		for _, key := range keys {
			h.Sum(rv.MapIndex(kv[key]).Bytes())
		}
		return true
	}
	return false
}

func structStrategy(h hash.Hash, rv reflect.Value) bool {
	if rv.Kind() == reflect.Struct {
		vtype := rv.Type()
		flen := vtype.NumField()
		fieldMap := make(map[string]reflect.Value, flen)

		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			ok, err := includeFieldPredicate(field, rv.Field(i))
			if err != nil && strings.Contains(err.Error(), "method:") {
				panic(err)
			}
			if !ok {
				continue
			}
			fieldMap[field.Name] = rv.Field(i)
		}

		keys := keysToSlice[string, reflect.Value](fieldMap)
		sort.Strings(keys)
		for _, k := range keys {
			v := fieldMap[k]
			h.Sum(v.Bytes())
		}
		return true
	}
	return false
}

func interfaceStrategy(h hash.Hash, rv reflect.Value) bool {
	if rv.Kind() == reflect.Interface {
		if !rv.CanInterface() {
			return false
		}
		strategies{
			primitiveStrategy,
			pointerStrategy,
			mapStrategy,
			structStrategy,
			interfaceStrategy,
			defaultStrategy,
		}.apply(h, reflect.ValueOf(rv.Interface()))
		return true
	}
	return false
}

func defaultStrategy(h hash.Hash, rv reflect.Value) bool {
	h.Sum(rv.Bytes())
	return true
}

type strategies []func(h hash.Hash, rv reflect.Value) bool

var identityStrategies = strategies{
	primitiveStrategy,
	pointerStrategy,
	mapStrategy,
	structStrategy,
	interfaceStrategy,
	defaultStrategy,
}

func (is strategies) apply(h hash.Hash, object any) {
	for _, strategy := range is {
		if strategy(h, reflect.ValueOf(object)) {
			break
		}
	}
}

func hashAny[T any](object T, h hash.Hash) []byte {
	identityStrategies.apply(h, object)
	return h.Sum(nil)
}

func HashIdentity[T any](t T) string {
	return hex.EncodeToString(hashAny(t, sha512.New()))
}
