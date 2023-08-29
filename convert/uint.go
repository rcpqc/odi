package convert

import (
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

var cvrsUint [types.MaxKinds]func(src reflect.Value) (uint64, error)

func init() {
	cvrsUint[reflect.Invalid] = cvrUintFromInvalid
	cvrsUint[reflect.Interface] = cvrUintFromInterface
	cvrsUint[reflect.Pointer] = cvrUintFromPointer
	cvrsUint[reflect.Bool] = cvrUintFromBool
	cvrsUint[reflect.Int] = cvrUintFromInt
	cvrsUint[reflect.Int8] = cvrUintFromInt
	cvrsUint[reflect.Int16] = cvrUintFromInt
	cvrsUint[reflect.Int32] = cvrUintFromInt
	cvrsUint[reflect.Int64] = cvrUintFromInt
	cvrsUint[reflect.Uint] = cvrUintFromUint
	cvrsUint[reflect.Uint8] = cvrUintFromUint
	cvrsUint[reflect.Uint16] = cvrUintFromUint
	cvrsUint[reflect.Uint32] = cvrUintFromUint
	cvrsUint[reflect.Uint64] = cvrUintFromUint
	cvrsUint[reflect.Uintptr] = cvrUintFromUint
	cvrsUint[reflect.Float32] = cvrUintFromFloat
	cvrsUint[reflect.Float64] = cvrUintFromFloat
	cvrsUint[reflect.String] = cvrUintFromString
}

func cvrUintFromInvalid(src reflect.Value) (uint64, error) {
	return 0, nil
}
func cvrUintFromInterface(src reflect.Value) (uint64, error) {
	if src.Type() == types.Any {
		return Uint(src.Elem())
	}
	return 0, errs.Newf("can't convert interface to integer")
}
func cvrUintFromPointer(src reflect.Value) (uint64, error) {
	if src.IsNil() {
		return 0, errs.Newf("can't nil pointer convert to integer")
	}
	return Uint(src.Elem())
}

func cvrUintFromBool(src reflect.Value) (uint64, error) {
	if src.Bool() {
		return 1, nil
	} else {
		return 0, nil
	}
}

func cvrUintFromInt(src reflect.Value) (uint64, error) {
	return uint64(src.Int()), nil
}

func cvrUintFromUint(src reflect.Value) (uint64, error) {
	return src.Uint(), nil
}

func cvrUintFromFloat(src reflect.Value) (uint64, error) {
	return uint64(src.Float()), nil
}

func cvrUintFromString(src reflect.Value) (uint64, error) {
	u64, err := strconv.ParseUint(src.String(), 10, 64)
	if err != nil {
		return 0, errs.Newf("string(%s) can't convert to integer", src.String())
	}
	return u64, nil
}

func Uint(src reflect.Value) (uint64, error) {
	if cvr := cvrsUint[src.Kind()]; cvr != nil {
		return cvr(src)
	}
	return 0, errs.Newf("can't convert kind(%v) to uint", src.Kind())
}
