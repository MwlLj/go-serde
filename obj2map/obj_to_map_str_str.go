package obj2map

import (
    "reflect"
    "errors"
    "fmt"
    "strconv"
)

var _ = fmt.Println

const (
    tag_field string = "field"
)

/*
** 传入一个字符串, 字符串中的 {var} 使用 结构体 tag 定义的字段的值来替换
*/
func Obj2MapStrStr(obj interface{}) (*map[string]string, error) {
    result := map[string]string{}
    err := Obj2MapStrStrWithCollect(obj, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func obj2MapStrStrStructInner(value reflect.Value, valueType reflect.Type, maps *map[string]string) {
    result := *maps
    fieldNum := value.NumField()
    for i := 0; i < fieldNum; i++ {
        field := value.Field(i)
        fieldType := valueType.Field(i)
        var fieldName string
        fieldTag := fieldType.Tag
        fieldName = fieldTag.Get(tag_field)
        if fieldName == "" {
            fieldName = fieldType.Name
        }
        fieldKind := field.Type().Kind()
        var fieldValue string
        switch fieldKind {
        case reflect.Ptr:
            fieldPtrType := field.Type().Elem()
            fieldPtrKind := fieldPtrType.Kind()
            fieldElem := field.Elem()
            switch fieldPtrKind {
            case reflect.String:
                fieldValue = fieldElem.String()
            case reflect.Int:
                i := fieldElem.Int()
                fieldValue = strconv.FormatInt(i, 10)
			case reflect.Int64:
                i := fieldElem.Int()
                fieldValue = strconv.FormatInt(i, 10)
            case reflect.Bool:
                b := fieldElem.Bool()
                fieldValue = strconv.FormatBool(b)
            case reflect.Float64:
                f := fieldElem.Float()
                fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
            }
        default:
            switch fieldKind {
            case reflect.String:
                fieldValue = field.String()
            case reflect.Int:
                i := field.Int()
                fieldValue = strconv.FormatInt(i, 10)
            case reflect.Int64:
                i := field.Int()
                fieldValue = strconv.FormatInt(i, 10)
            case reflect.Bool:
                b := field.Bool()
                fieldValue = strconv.FormatBool(b)
            case reflect.Float64:
                f := field.Float()
                fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
            }
        }
        result[fieldName] = fieldValue
    }
}

func calcFieldValue(field reflect.Value, fieldValue *string) {
    fieldKind := field.Type().Kind()
    switch fieldKind {
    case reflect.Ptr:
        fieldPtrType := field.Type().Elem()
        fieldPtrKind := fieldPtrType.Kind()
        fieldElem := field.Elem()
        switch fieldPtrKind {
        case reflect.String:
            *fieldValue = fieldElem.String()
        case reflect.Int:
            i := fieldElem.Int()
            *fieldValue = strconv.FormatInt(i, 10)
        case reflect.Int64:
            i := fieldElem.Int()
            *fieldValue = strconv.FormatInt(i, 10)
        case reflect.Bool:
            b := fieldElem.Bool()
            *fieldValue = strconv.FormatBool(b)
        case reflect.Float64:
            f := fieldElem.Float()
            *fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
        case reflect.Interface:
            *fieldValue = fieldElem.String()
        }
    default:
        switch fieldKind {
        case reflect.String:
            *fieldValue = field.String()
        case reflect.Int:
            i := field.Int()
            *fieldValue = strconv.FormatInt(i, 10)
        case reflect.Int64:
            i := field.Int()
            *fieldValue = strconv.FormatInt(i, 10)
        case reflect.Bool:
            b := field.Bool()
            *fieldValue = strconv.FormatBool(b)
        case reflect.Float64:
            f := field.Float()
            *fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
        case reflect.Interface:
            calcFieldValue(field.Elem(), fieldValue)
        }
    }
}

func obj2MapStrStrMapInner(value reflect.Value, valueType reflect.Type, maps *map[string]string) {
    result := *maps
    for _, key := range value.MapKeys() {
        // field := value.Field(i)
        // fieldType := valueType.Field(i)
        field := value.MapIndex(key)
        fieldName := key.String()
        // fieldKind := field.Type().Kind()
        var fieldValue string
        calcFieldValue(field, &fieldValue)
        result[fieldName] = fieldValue
    }
}

func Obj2MapStrStrWithCollect(obj interface{}, maps *map[string]string) error {
    if maps == nil {
        return errors.New("param is nil")
    }
    var valueType reflect.Type
    value := reflect.ValueOf(obj)
    outType := reflect.TypeOf(obj)
    outputKind := outType.Kind()
    switch outputKind {
    case reflect.Ptr:
        valueType = outType.Elem()
        value = value.Elem()
    default:
        valueType = reflect.TypeOf(obj)
    }
    /*
    ** valueType: 指针指向的类型
    */
    switch valueType.Kind() {
    case reflect.Struct:
        obj2MapStrStrStructInner(value, valueType, maps)
    case reflect.Interface:
        obj2MapStrStrStructInner(value, valueType, maps)
    case reflect.Map:
        obj2MapStrStrMapInner(value, valueType, maps)
    default:
        return errors.New("is not parse")
    }
    return nil
}
