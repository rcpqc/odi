package resolve

import (
	"context"
	"reflect"
	"strconv"

	"github.com/rcpqc/odi/errs"
)

func injectBool(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() == reflect.Bool {
		dst.SetBool(src.Bool())
		return nil
	}
	if src.CanInt() {
		dst.SetBool(src.Int() != 0)
		return nil
	}
	if src.CanUint() {
		dst.SetBool(src.Uint() != 0)
		return nil
	}
	if src.CanFloat() {
		dst.SetBool(src.Float() != 0)
		return nil
	}
	if src.Kind() == reflect.String {
		b, err := strconv.ParseBool(src.String())
		if err != nil {
			return err
		}
		dst.SetBool(b)
		return nil
	}
	return errs.Newf("can't inject type(%v) expect type(%v)", src.Type(), dst.Type())
}

func injectInt(ctx context.Context, dst, src reflect.Value) error {
	if src.CanInt() {
		dst.SetInt(src.Int())
		return nil
	}
	if src.CanUint() {
		dst.SetInt(int64(src.Uint()))
		return nil
	}
	if src.CanFloat() {
		dst.SetInt(int64(src.Float()))
		return nil
	}
	if src.Kind() == reflect.String {
		i64, err := strconv.ParseInt(src.String(), 10, 64)
		if err != nil {
			return errs.Newf("string(%s) can't convert to integer", src.String())
		}
		dst.SetInt(i64)
		return nil
	}
	if src.Kind() == reflect.Bool {
		if src.Bool() {
			dst.SetInt(1)
		} else {
			dst.SetInt(0)
		}
		return nil
	}
	return errs.Newf("can't inject type(%v) expect type(%v)", src.Type(), dst.Type())
}

func injectUint(ctx context.Context, dst, src reflect.Value) error {
	if src.CanInt() {
		dst.SetUint(uint64(src.Int()))
		return nil
	}
	if src.CanUint() {
		dst.SetUint(src.Uint())
		return nil
	}
	if src.CanFloat() {
		dst.SetUint(uint64(src.Float()))
		return nil
	}
	if src.Kind() == reflect.String {
		u64, _ := strconv.ParseUint(src.String(), 10, 64)
		dst.SetUint(u64)
		return nil
	}
	if src.Kind() == reflect.Bool {
		if src.Bool() {
			dst.SetUint(1)
		} else {
			dst.SetUint(0)
		}
		return nil
	}
	return errs.Newf("can't inject type(%v) expect type(%v)", src.Type(), dst.Type())
}

func injectFloat(ctx context.Context, dst, src reflect.Value) error {
	if src.CanFloat() {
		dst.SetFloat(src.Float())
		return nil
	}
	if src.CanInt() {
		dst.SetFloat(float64(src.Int()))
		return nil
	}
	if src.CanUint() {
		dst.SetFloat(float64(src.Uint()))
		return nil
	}
	if src.Kind() == reflect.String {
		f64, _ := strconv.ParseFloat(src.String(), 64)
		dst.SetFloat(f64)
		return nil
	}
	return errs.Newf("can't inject type(%v) expect type(%v)", src.Type(), dst.Type())
}

func injectString(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() == reflect.String {
		dst.SetString(src.String())
		return nil
	}
	if src.CanInt() {
		dst.SetString(strconv.FormatInt(src.Int(), 10))
		return nil
	}
	if src.CanUint() {
		dst.SetString(strconv.FormatUint(src.Uint(), 10))
		return nil
	}
	if src.CanFloat() {
		dst.SetString(strconv.FormatFloat(src.Float(), 'f', 6, 64))
		return nil
	}
	if src.Kind() == reflect.Bool {
		if src.Bool() {
			dst.SetString("true")
		} else {
			dst.SetString("false")
		}
		return nil
	}
	return errs.Newf("can't inject type(%v) expect type(%v)", src.Type(), dst.Type())
}
