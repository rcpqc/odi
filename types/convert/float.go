package convert

import (
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

var cvrsFloat [types.MaxKinds]func(src reflect.Value) (float64, error)

func init() {
	cvrsFloat[reflect.Invalid] = cvrFloatFromInvalid
	cvrsFloat[reflect.Pointer] = cvrFloatFromPointer
	cvrsFloat[reflect.Bool] = cvrFloatFromBool
	cvrsFloat[reflect.Int] = cvrFloatFromInt
	cvrsFloat[reflect.Int8] = cvrFloatFromInt
	cvrsFloat[reflect.Int16] = cvrFloatFromInt
	cvrsFloat[reflect.Int32] = cvrFloatFromInt
	cvrsFloat[reflect.Int64] = cvrFloatFromInt
	cvrsFloat[reflect.Uint] = cvrFloatFromUint
	cvrsFloat[reflect.Uint8] = cvrFloatFromUint
	cvrsFloat[reflect.Uint16] = cvrFloatFromUint
	cvrsFloat[reflect.Uint32] = cvrFloatFromUint
	cvrsFloat[reflect.Uint64] = cvrFloatFromUint
	cvrsFloat[reflect.Uintptr] = cvrFloatFromUint
	cvrsFloat[reflect.Float32] = cvrFloatFromFloat
	cvrsFloat[reflect.Float64] = cvrFloatFromFloat
	cvrsFloat[reflect.String] = cvrFloatFromString
}

func cvrFloatFromInvalid(src reflect.Value) (float64, error) {
	return 0.0, nil
}

func cvrFloatFromPointer(src reflect.Value) (float64, error) {
	if src.IsNil() {
		return 0, errs.Newf("can't convert nil pointer to float")
	}
	return Float(src.Elem())
}

func cvrFloatFromBool(src reflect.Value) (float64, error) {
	if src.Bool() {
		return 1.0, nil
	} else {
		return 0.0, nil
	}
}

func cvrFloatFromInt(src reflect.Value) (float64, error) {
	return float64(src.Int()), nil
}

func cvrFloatFromUint(src reflect.Value) (float64, error) {
	return float64(src.Uint()), nil
}

func cvrFloatFromFloat(src reflect.Value) (float64, error) {
	return src.Float(), nil
}

func cvrFloatFromString(src reflect.Value) (float64, error) {
	f64, err := strconv.ParseFloat(src.String(), 64)
	if err != nil {
		return 0.0, errs.Newf("can't convert string(%s) to float", src.String())
	}
	return f64, nil
}

func Float(src reflect.Value) (float64, error) {
	if cvr := cvrsFloat[src.Kind()]; cvr != nil {
		return cvr(src)
	}
	return 0, errs.Newf("can't convert kind(%v) to float", src.Kind())
}
