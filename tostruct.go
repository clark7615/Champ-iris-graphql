package graphql

import (
	"reflect"
	"strings"
	"time"
)

//ToStruct Map[String]interface to Struct
func ToStruct(in interface{}, out interface{}) {
	v := reflect.ValueOf(out)
	if v.IsValid() == false {
		panic("not valid")
	}
	//找到最后指向的值，或者空指针，空指针是需要进行初始化的
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}

	tv := v
	if tv.Kind() == reflect.Ptr && tv.CanSet() {
		//对空指针进行初始化，暂时用临时变量保存
		tv.Set(reflect.New(tv.Type().Elem()))
		tv = tv.Elem()
	}

	if tv.Kind() != reflect.Struct {
		panic("not struct")
	}

	if assign(tv, in, "json") { //赋值成功，将临时变量赋给原值
		if v.Kind() == reflect.Ptr {
			v.Set(tv.Addr())
		} else {
			v.Set(tv)
		}
	}
}

func assign(dstVal reflect.Value, src interface{}, tagName string) bool {
	sv := reflect.ValueOf(src)
	if !dstVal.IsValid() || !sv.IsValid() {
		return false
	}

	if dstVal.Kind() == reflect.Ptr {
		//初始化空指针
		if dstVal.IsNil() && dstVal.CanSet() {
			dstVal.Set(reflect.New(dstVal.Type().Elem()))
		}
		dstVal = dstVal.Elem()
	}

	// 判断可否赋值，小写字母开头的字段、常量等不可赋值
	if !dstVal.CanSet() {
		return false
	}

	switch dstVal.Kind() {
	case reflect.Bool:
		dstVal.Set(reflect.ValueOf(src))
		return true
	case reflect.Float32, reflect.Float64:
		return coerceFloat(dstVal, src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return coerceInt(dstVal, src)
	case reflect.String:
		dstVal.Set(reflect.ValueOf(src))
		return true
	case reflect.Ptr, reflect.Slice, reflect.Map:
		dstVal.Set(reflect.Zero(dstVal.Type()))
		return true
	case reflect.Struct:
		if dstVal.Type() == reflect.TypeOf(time.Time{}) {
			parse, _ := time.ParseInLocation("2006/01/02", reflect.ValueOf(src).String(), time.Local)
			dstVal.Set(reflect.ValueOf(parse))
			return true
		}

		if sv.Kind() != reflect.Map || sv.Type().Key().Kind() != reflect.String {
			return false
		}

		success := false
		for i := 0; i < dstVal.NumField(); i++ {
			fv := dstVal.Field(i)
			if fv.IsValid() == false || fv.CanSet() == false {
				continue
			}

			ft := dstVal.Type().Field(i)
			name := ft.Name
			strs := strings.Split(ft.Tag.Get(tagName), ",")
			if strs[0] == "-" { //处理ignore的标志
				continue
			}

			if len(strs[0]) > 0 {
				name = strs[0]
			}
			fsv := sv.MapIndex(reflect.ValueOf(name))
			if fsv.IsValid() {
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					pv := reflect.New(fv.Type().Elem())
					if assign(pv, fsv.Interface(), tagName) {
						fv.Set(pv)
						success = true
					}
				} else {
					if assign(fv, fsv.Interface(), tagName) {
						success = true
					}
				}
			} else if ft.Anonymous {
				//尝试对匿名字段进行递归赋值，跟JSON的处理原则保持一致
				if fv.Kind() == reflect.Ptr && fv.IsNil() {
					pv := reflect.New(fv.Type().Elem())
					if assign(pv, src, tagName) {
						fv.Set(pv)
						success = true
					}
				} else {
					if assign(fv, src, tagName) {
						success = true
					}
				}
			}
		}
		return success
	default:
		return false
	}
}

func coerceInt(val reflect.Value, src interface{}) bool {
	switch val.Kind() {
	case reflect.Int:
		val.Set(reflect.ValueOf(src))
	case reflect.Int8:
		i := int8(src.(int))
		val.Set(reflect.ValueOf(i))
	case reflect.Int16:
		i := int16(src.(int))
		val.Set(reflect.ValueOf(i))
	case reflect.Int32:
		i := int32(src.(int))
		val.Set(reflect.ValueOf(i))
	case reflect.Int64:
		i := int64(src.(int))
		val.Set(reflect.ValueOf(i))
	default:
		return false
	}
	return true
}

func coerceFloat(val reflect.Value, src interface{}) bool {
	switch val.Kind() {
	case reflect.Float64:
		val.Set(reflect.ValueOf(src))
	case reflect.Float32:
		i := float32(src.(float64))
		val.Set(reflect.ValueOf(i))
	default:
		return false
	}
	return true
}
