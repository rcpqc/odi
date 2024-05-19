package resolve

import (
	"context"
	"reflect"

	"github.com/rcpqc/odi/errs"
	"github.com/rcpqc/odi/types"
	"github.com/rcpqc/odi/types/convert"
)

// sourceConvert convert map key to string
func sourceConvert(src reflect.Value) reflect.Value {
	m := map[string]interface{}{}
	for iter := src.MapRange(); iter.Next(); {
		key := iter.Key()
		if key.Type() == types.Any {
			key = key.Elem()
		}
		if s, err := convert.String(key); err == nil {
			m[s] = iter.Value().Interface()
		}
	}
	return reflect.ValueOf(m)
}

// sourceExclude exclude some keys in source
func sourceExclude(src reflect.Value, excludes map[any]struct{}) reflect.Value {
	dst := reflect.MakeMap(src.Type())
	for iter := src.MapRange(); iter.Next(); {
		key, val := iter.Key(), iter.Value()
		if _, ok := excludes[key.Interface()]; ok {
			continue
		}
		dst.SetMapIndex(key, val)
	}
	return dst
}

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if !src.IsValid() {
		return nil
	}
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	if ctxGetStructFieldNameCompatibility(ctx) {
		src = sourceConvert(src)
	}
	tProfile := ctxGetProfileFactory(ctx).GetProfile(dst.Type())
	for _, field := range tProfile.Fields {
		if field.Error != nil {
			return errs.New(field.Error).Prefix(field.Router)
		}
		vfield := dst.FieldByIndex(field.Index)
		if !vfield.CanSet() {
			continue
		}
		if field.InlineMap {
			if err := injectMap(ctx, vfield, sourceExclude(src, tProfile.Names)); err != nil {
				return errs.New(err).Prefix(field.Router)
			}
			continue
		}
		if field.InlineIface {
			if err := injectInterface(ctx, vfield, sourceExclude(src, tProfile.Names)); err != nil {
				return errs.New(err).Prefix(field.Router)
			}
			continue
		}
		key := reflect.ValueOf(field.Name)
		val := src.MapIndex(key)
		if field.Required && !val.IsValid() {
			return errs.Newf("required but not assigned").Prefix(field.Router)
		}
		if !val.IsValid() {
			continue
		}
		if err := inject(ctx, vfield, val); err != nil {
			return errs.New(err).Prefix(field.Router)
		}
	}
	return nil
}
