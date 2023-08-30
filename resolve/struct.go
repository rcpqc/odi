package resolve

import (
	"context"
	"fmt"
	"reflect"

	"github.com/rcpqc/odi/convert"
	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
)

// convertSource convert map key to string
func convertSource(src reflect.Value) reflect.Value {
	m := map[string]interface{}{}
	for iter := src.MapRange(); iter.Next(); {
		if s, err := convert.String(iter.Key()); err == nil {
			m[s] = iter.Value().Interface()
		}
	}
	return reflect.ValueOf(m)
}

func injectStructInlineMap(ctx context.Context, dst, src reflect.Value, excludes map[string]struct{}) error {
	if dst.IsNil() {
		dst.Set(reflect.MakeMap(dst.Type()))
	}
	iter := src.MapRange()
	for iter.Next() {
		if _, ok := excludes[iter.Key().String()]; ok {
			continue
		}
		srcKey, srcVal := iter.Key(), iter.Value()
		dstKey, dstVal := reflect.New(dst.Type().Key()).Elem(), reflect.New(dst.Type().Elem()).Elem()
		if err := inject(ctx, dstKey, srcKey); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%s]", dstKey.String()))
		}
		if err := inject(ctx, dstVal, srcVal); err != nil {
			return errs.New(err).Prefix(fmt.Sprintf("[%s]", dstKey.String()))
		}
		dst.SetMapIndex(dstKey, dstVal)
	}
	return nil
}

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	if ctxGetStructFieldNameCompatibility(ctx) {
		src = convertSource(src)
	}
	tProfile := types.GetProfile(dst.Type(), ctxGetTagKey(ctx))
	for _, field := range tProfile.Fields {
		if field.Error != nil {
			return errs.New(field.Error).Prefix("." + field.Name)
		}
		vfield := dst.FieldByIndex(field.Index)
		if !vfield.CanSet() {
			continue
		}
		if field.InlineMap {
			if err := injectStructInlineMap(ctx, vfield, src, tProfile.Names); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
			continue
		}
		key := reflect.ValueOf(field.Name)
		val := src.MapIndex(key)
		if !val.IsValid() {
			continue
		}
		if err := inject(ctx, vfield, val); err != nil {
			return errs.New(err).Prefix("." + field.Name)
		}
	}
	return nil
}
