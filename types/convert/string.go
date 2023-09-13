package convert

import (
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

var cvrsString [types.MaxKinds]func(src reflect.Value) (string, error)

func init() {
	cvrsString[reflect.Invalid] = cvrStringFromInvalid
	cvrsString[reflect.Pointer] = cvrStringFromPointer
	cvrsString[reflect.Bool] = cvrStringFromBool
	cvrsString[reflect.Int] = cvrStringFromInt
	cvrsString[reflect.Int8] = cvrStringFromInt
	cvrsString[reflect.Int16] = cvrStringFromInt
	cvrsString[reflect.Int32] = cvrStringFromInt
	cvrsString[reflect.Int64] = cvrStringFromInt
	cvrsString[reflect.Uint] = cvrStringFromUint
	cvrsString[reflect.Uint8] = cvrStringFromUint
	cvrsString[reflect.Uint16] = cvrStringFromUint
	cvrsString[reflect.Uint32] = cvrStringFromUint
	cvrsString[reflect.Uint64] = cvrStringFromUint
	cvrsString[reflect.Uintptr] = cvrStringFromUint
	cvrsString[reflect.Float32] = cvrStringFromFloat
	cvrsString[reflect.Float64] = cvrStringFromFloat
	cvrsString[reflect.String] = cvrStringFromString
}

func cvrStringFromInvalid(src reflect.Value) (string, error) {
	return "", nil
}

func cvrStringFromPointer(src reflect.Value) (string, error) {
	if src.IsNil() {
		return "", errs.Newf("can't convert nil pointer to string")
	}
	return String(src.Elem())
}

func cvrStringFromBool(src reflect.Value) (string, error) {
	if src.Bool() {
		return "true", nil
	} else {
		return "false", nil
	}
}

func cvrStringFromInt(src reflect.Value) (string, error) {
	return strconv.FormatInt(src.Int(), 10), nil
}

func cvrStringFromUint(src reflect.Value) (string, error) {
	return strconv.FormatUint(src.Uint(), 10), nil
}

func cvrStringFromFloat(src reflect.Value) (string, error) {
	return strconv.FormatFloat(src.Float(), 'f', -1, 64), nil
}

func cvrStringFromString(src reflect.Value) (string, error) {
	return src.String(), nil
}

func String(src reflect.Value) (string, error) {
	if cvr := cvrsString[src.Kind()]; cvr != nil {
		return cvr(src)
	}
	return "", errs.Newf("can't convert kind(%v) to string", src.Kind())
}
