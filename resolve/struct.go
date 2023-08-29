package resolve

import (
	"context"
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

func injectStructInternal(ctx context.Context, dst, src reflect.Value) error {
	for _, field := range types.NewProfile(dst.Type(), ctxGetTagKey(ctx)).GetFields() {
		vfield := dst.Field(field.Index)
		if !vfield.CanSet() {
			continue
		}
		// direct injection of non-inline fields
		if !field.Inline {
			key := reflect.ValueOf(field.Name)
			val := src.MapIndex(key)
			if !val.IsValid() {
				continue
			}
			if err := inject(ctx, vfield, val); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
			// delete injected fields to prevent repeated injection of inline fields
			src.SetMapIndex(key, reflect.Value{})
			continue
		}
		// inline fields are injected according to type (map[string]any/struct/*struct)
		if field.Type.Kind() == reflect.Map && field.Type.Key().Kind() == reflect.String {
			if err := inject(ctx, vfield, src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else if field.Type.Kind() == reflect.Struct {
			if err := injectStructInternal(ctx, vfield, src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			if vfield.IsNil() {
				vfield.Set(reflect.New(vfield.Type().Elem()))
			}
			if err := injectStructInternal(ctx, vfield.Elem(), src); err != nil {
				return errs.New(err).Prefix("." + field.Name)
			}
		} else {
			return errs.Newf("illegal inline type(%v) expect struct or map[string]", field.Type).Prefix("." + field.Name)
		}
		continue
	}
	return nil
}

func injectStruct(ctx context.Context, dst, src reflect.Value) error {
	if src.Kind() != reflect.Map {
		return errs.Newf("expect map but %v", src.Kind())
	}
	return injectStructInternal(ctx, dst, convertSource(src))
}
