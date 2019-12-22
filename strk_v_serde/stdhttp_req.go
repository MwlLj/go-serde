package kv_serde

import (
    "reflect"
    "net/http"
    "errors"
    "strconv"
    "fmt"
    "encoding/json"
)

var _ = fmt.Println

const (
    tag_field string = "field"
    tag_type string = "type"
    tag_type_json string = "json"
)

type kv interface {
    init(req *http.Request) error
    get(key *string) string
}

/*
** 反序列化 string类型的key, 任意类型的 value
*/
func deserde(input kv, req *http.Request, output interface{}) error {
    outValue := reflect.ValueOf(output)
    if outValue.IsNil() {
        return nil
    }
    outType := reflect.TypeOf(output)
    outputKind := outType.Kind()
    var valueType reflect.Type
    switch outputKind {
    case reflect.Ptr:
        valueType = outType.Elem()
    default:
        return errors.New("is not ptr")
    }
    value := outValue.Elem()
    /*
    ** valueType: 指针指向的类型
    */
    switch valueType.Kind() {
    case reflect.Struct:
    default:
        return errors.New("is not struct")
    }
    fieldNum := value.NumField()
    if err := input.init(req); err != nil {
        return err
    }
    for i := 0; i < fieldNum; i++ {
        field := value.Field(i)
        fieldType := valueType.Field(i)
        var fieldName string
        fieldTag := fieldType.Tag
        fieldName = fieldTag.Get(tag_field)
        if fieldName == "" {
            fieldName = fieldType.Name
        }
        /*
        ** 判断字段类型
        */
        fieldKind := field.Type().Kind()
        switch fieldKind {
        case reflect.Ptr:
            /*
            ** 字段为指针
            */
            fieldPtrType := field.Type().Elem()
            fieldPtrKind := fieldPtrType.Kind()
            switch fieldPtrKind {
            case reflect.String:
                v := input.get(&fieldName)
                if v != "" {
                    fieldPtrValue := reflect.New(fieldPtrType)
                    fieldPtrValueElem := fieldPtrValue.Elem()
                    fieldPtrValueElem.SetString(v)
                    field.Set(fieldPtrValue)
                }
            case reflect.Int:
                v := input.get(&fieldName)
                if v != "" {
                    fieldPtrValue := reflect.New(fieldPtrType)
                    fieldPtrValueElem := fieldPtrValue.Elem()
                    if vi, err := strconv.ParseInt(v, 10, 64); err == nil {
                        fieldPtrValueElem.SetInt(vi)
                    } else {
                        fieldPtrValueElem.SetInt(0)
                    }
                    field.Set(fieldPtrValue)
                }
            case reflect.Bool:
                v := input.get(&fieldName)
                if v != "" {
                    fieldPtrValue := reflect.New(fieldPtrType)
                    fieldPtrValueElem := fieldPtrValue.Elem()
                    if vi, err := strconv.ParseBool(v); err == nil {
                        fieldPtrValueElem.SetBool(vi)
                    } else {
                        fieldPtrValueElem.SetBool(false)
                    }
                    field.Set(fieldPtrValue)
                }
            case reflect.Float64:
                v := input.get(&fieldName)
                if v != "" {
                    fieldPtrValue := reflect.New(fieldPtrType)
                    fieldPtrValueElem := fieldPtrValue.Elem()
                    if vi, err := strconv.ParseFloat(v, 64); err == nil {
                        fieldPtrValueElem.SetFloat(vi)
                    } else {
                        fieldPtrValueElem.SetFloat(0.0)
                    }
                    field.Set(fieldPtrValue)
                }
            case reflect.Struct:
                v := input.get(&fieldName)
                if v != "" {
                    fieldPtrValue := reflect.New(fieldPtrType)
                    fieldPtrValueElem := fieldPtrValue.Elem()
                    fieldValue := reflect.New(fieldPtrType)
                    json.Unmarshal([]byte(v), fieldValue.Interface())
                    fieldPtrValueElem.Set(fieldValue.Elem())
                    field.Set(fieldPtrValue)
                }
            }
        default:
            /*
            ** 字段不为指针
            */
            switch fieldKind {
            case reflect.String:
                v := input.get(&fieldName)
                field.SetString(v)
            case reflect.Int:
                v := input.get(&fieldName)
                if vi, err := strconv.ParseInt(v, 10, 64); err == nil {
                    field.SetInt(vi)
                } else {
                    field.SetInt(0)
                }
            case reflect.Bool:
                v := input.get(&fieldName)
                if vi, err := strconv.ParseBool(v); err == nil {
                    field.SetBool(vi)
                } else {
                    field.SetBool(false)
                }
            case reflect.Float64:
                v := input.get(&fieldName)
                if vi, err := strconv.ParseFloat(v, 64); err == nil {
                    field.SetFloat(vi)
                } else {
                    field.SetFloat(0.0)
                }
            case reflect.Struct:
                v := input.get(&fieldName)
                fieldValue := reflect.New(field.Type())
                json.Unmarshal([]byte(v), fieldValue.Interface())
                field.Set(fieldValue.Elem())
                // tmp := field.Interface()
                // json.Unmarshal([]byte(v), tmp)
                // field.Set(reflect.ValueOf(tmp))
            }
        }
    }
    return nil
}
