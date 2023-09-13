package convert

import (
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

var cvrsInt [types.MaxKinds]func(src reflect.Value) (int64, error)

func init() {
	cvrsInt[reflect.Pointer] = cvrIntFromPointer
	cvrsInt[reflect.Bool] = cvrIntFromBool
	cvrsInt[reflect.Int] = cvrIntFromInt
	cvrsInt[reflect.Int8] = cvrIntFromInt
	cvrsInt[reflect.Int16] = cvrIntFromInt
	cvrsInt[reflect.Int32] = cvrIntFromInt
	cvrsInt[reflect.Int64] = cvrIntFromInt
	cvrsInt[reflect.Uint] = cvrIntFromUint
	cvrsInt[reflect.Uint8] = cvrIntFromUint
	cvrsInt[reflect.Uint16] = cvrIntFromUint
	cvrsInt[reflect.Uint32] = cvrIntFromUint
	cvrsInt[reflect.Uint64] = cvrIntFromUint
	cvrsInt[reflect.Uintptr] = cvrIntFromUint
	cvrsInt[reflect.Float32] = cvrIntFromFloat
	cvrsInt[reflect.Float64] = cvrIntFromFloat
	cvrsInt[reflect.String] = cvrIntFromString
}

func cvrIntFromPointer(src reflect.Value) (int64, error) {
	if src.IsNil() {
		return 0, errs.Newf("can't convert nil pointer to int")
	}
	return Int(src.Elem())
}

func cvrIntFromBool(src reflect.Value) (int64, error) {
	if src.Bool() {
		return 1, nil
	} else {
		return 0, nil
	}
}

func cvrIntFromInt(src reflect.Value) (int64, error) {
	return src.Int(), nil
}

func cvrIntFromUint(src reflect.Value) (int64, error) {
	return int64(src.Uint()), nil
}

func cvrIntFromFloat(src reflect.Value) (int64, error) {
	return int64(src.Float()), nil
}

func cvrIntFromString(src reflect.Value) (int64, error) {
	i64, err := strconv.ParseInt(src.String(), 10, 64)
	if err != nil {
		return 0, errs.Newf("can't convert string(%s) to int", src.String())
	}
	return i64, nil
}

func Int(src reflect.Value) (int64, error) {
	if cvr := cvrsInt[src.Kind()]; cvr != nil {
		return cvr(src)
	}
	return 0, errs.Newf("can't convert kind(%v) to int", src.Kind())
}
