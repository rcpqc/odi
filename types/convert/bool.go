package convert

import (
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

var cvrsBool [types.MaxKinds]func(src reflect.Value) (bool, error)

func init() {
	cvrsBool[reflect.Invalid] = cvrBoolFromInvalid
	cvrsBool[reflect.Pointer] = cvrBoolFromPointer
	cvrsBool[reflect.Bool] = cvrBoolFromBool
	cvrsBool[reflect.Int] = cvrBoolFromInt
	cvrsBool[reflect.Int8] = cvrBoolFromInt
	cvrsBool[reflect.Int16] = cvrBoolFromInt
	cvrsBool[reflect.Int32] = cvrBoolFromInt
	cvrsBool[reflect.Int64] = cvrBoolFromInt
	cvrsBool[reflect.Uint] = cvrBoolFromUint
	cvrsBool[reflect.Uint8] = cvrBoolFromUint
	cvrsBool[reflect.Uint16] = cvrBoolFromUint
	cvrsBool[reflect.Uint32] = cvrBoolFromUint
	cvrsBool[reflect.Uint64] = cvrBoolFromUint
	cvrsBool[reflect.Uintptr] = cvrBoolFromUint
	cvrsBool[reflect.Float32] = cvrBoolFromFloat
	cvrsBool[reflect.Float64] = cvrBoolFromFloat
	cvrsBool[reflect.String] = cvrBoolFromString
}

func cvrBoolFromInvalid(src reflect.Value) (bool, error) {
	return false, nil
}

func cvrBoolFromPointer(src reflect.Value) (bool, error) {
	if src.IsNil() {
		return false, errs.Newf("can't convert nil pointer to bool")
	}
	return Bool(src.Elem())
}

func cvrBoolFromBool(src reflect.Value) (bool, error) {
	return src.Bool(), nil
}

func cvrBoolFromInt(src reflect.Value) (bool, error) {
	return src.Int() != 0, nil
}

func cvrBoolFromUint(src reflect.Value) (bool, error) {
	return src.Uint() != 0, nil
}

func cvrBoolFromFloat(src reflect.Value) (bool, error) {
	return src.Float() != 0, nil
}

func cvrBoolFromString(src reflect.Value) (bool, error) {
	b, err := strconv.ParseBool(src.String())
	if err != nil {
		return false, errs.Newf("can't convert string(%s) to bool", src.String())
	}
	return b, nil
}

func Bool(src reflect.Value) (bool, error) {
	if cvr := cvrsBool[src.Kind()]; cvr != nil {
		return cvr(src)
	}
	return false, errs.Newf("can't convert kind(%v) to bool", src.Kind())
}
