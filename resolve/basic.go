package resolve

import (
	"context"
	"reflect"

	"github.com/rcpqc/odi/types/convert"
)

func injectBool(ctx context.Context, dst, src reflect.Value) error {
	b, err := convert.Bool(src)
	if err != nil {
		return err
	}
	dst.SetBool(b)
	return nil
}

func injectInt(ctx context.Context, dst, src reflect.Value) error {
	i, err := convert.Int(src)
	if err != nil {
		return err
	}
	dst.SetInt(i)
	return nil
}

func injectUint(ctx context.Context, dst, src reflect.Value) error {
	u, err := convert.Uint(src)
	if err != nil {
		return err
	}
	dst.SetUint(u)
	return nil
}

func injectFloat(ctx context.Context, dst, src reflect.Value) error {
	f, err := convert.Float(src)
	if err != nil {
		return err
	}
	dst.SetFloat(f)
	return nil
}

func injectString(ctx context.Context, dst, src reflect.Value) error {
	s, err := convert.String(src)
	if err != nil {
		return err
	}
	dst.SetString(s)
	return nil
}
