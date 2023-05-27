package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/exp/maps"
)

type tagError string

func (e tagError) Error() string {
	return "incorrect tag " + string(e)
}

type structFieldFilter func(f reflect.StructField, v reflect.Value) (bool, error)

func filterField(f reflect.StructField, v reflect.Value) (bool, error) {
	if str := f.Tag.Get("hash"); str != "" {
		if str == "-" {
			return false, nil
		}
		for _, tag := range strings.Split(str, " ") {
			args := strings.Split(strings.TrimSpace(tag), ":")
			if len(args) != 2 {
				return false, tagError(tag)
			}
			switch args[0] {
			case "method":
				property, found := f.Type.MethodByName(strings.TrimSpace(args[1]))
				if !found || property.Type.NumOut() != 1 {
					return false, tagError(tag)
				}
				v = property.Func.Call([]reflect.Value{v})[0]
			}
		}
	}
	return true, nil
}

func sumValue(hasher hash.Hash, val reflect.Value, fltr structFieldFilter) {
	switch val.Kind() {
	case reflect.String:
		hasher.Sum(val.Bytes())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		hasher.Sum(val.Bytes())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		hasher.Sum(val.Bytes())
	case reflect.Float32, reflect.Float64:
		hasher.Sum(val.Bytes())
	case reflect.Bool:
		hasher.Sum(val.Bytes())
	case reflect.Ptr:
		if !val.IsNil() || val.Type().Elem().Kind() == reflect.Struct {
			sumValue(hasher, reflect.Indirect(val), fltr)
		} else {
			sumValue(hasher, reflect.Zero(val.Type().Elem()), fltr)
		}
	case reflect.Array, reflect.Slice:
		hasher.Sum(val.Bytes())
	case reflect.Map:
		mk := val.MapKeys()
		kv := make(map[string]reflect.Value, len(mk))
		for _, k := range mk {
			kv[k.String()] = k
		}
		keys := maps.Keys(kv)
		sort.Strings(keys)
		for _, key := range keys {
			hasher.Sum(val.MapIndex(kv[key]).Bytes())
		}
	case reflect.Struct:
		vtype := val.Type()
		flen := vtype.NumField()
		fieldMap := make(map[string]reflect.Value, flen)
		// Get all fields
		for i := 0; i < flen; i++ {
			field := vtype.Field(i)
			if fltr != nil {
				ok, err := fltr(field, val.Field(i))
				if err != nil && strings.Contains(err.Error(), "method:") {
					panic(err)
				}
				if !ok {
					continue
				}
			}
			fieldMap[field.Name] = val.Field(i)
		}

		keys := maps.Keys(fieldMap)
		sort.Strings(keys)
		for _, k := range keys {
			v := fieldMap[k]
			hasher.Sum(v.Bytes())
		}
	case reflect.Interface:
		if !val.CanInterface() {
			return
		}
		sumValue(hasher, reflect.ValueOf(val.Interface()), fltr)
	default:
		hasher.Sum(val.Bytes())
	}
}

func hashStruct[T any](object T, hasher hash.Hash) []byte {
	sumValue(hasher, reflect.ValueOf(object),
		func(f reflect.StructField, v reflect.Value) (bool, error) {
			return filterField(f, v)
		})

	return hasher.Sum(nil)
}

func HashStructIdentity[T any](t T) string {
	return hex.EncodeToString(hashStruct(t, sha256.New()))
}
