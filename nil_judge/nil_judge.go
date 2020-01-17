package nil_judge

import (
    "reflect"
    "fmt"
)

var _ = fmt.Println

/*
** 如果是结构体, 判断结构体中的所有字段是否都为空
*/
type CNilJudge struct {
}

/*
** 是否全是 nil
**  全是nil => true
**  有部分不是nil => false
*/
func (self *CNilJudge) IsAllNil(obj interface{}) bool {
    if obj == nil {
        return true
    }
    objType := reflect.TypeOf(obj)
    objValue := reflect.ValueOf(obj)
    objKind := objType.Kind()
    switch objKind {
    case reflect.Ptr:
        if objValue.IsNil() {
            return true
        }
        objElemValue := objValue.Elem()
        objElemType := objElemValue.Type()
        objElemKind := objElemType.Kind()
        switch objElemKind {
        case reflect.Struct:
            /*
            ** 判断结构体中的所有字段是否都为空
            */
            isAllNil := true
            num := objElemValue.NumField()
            for i := 0; i < num; i++ {
                field := objElemValue.Field(i)
                switch field.Type().Kind() {
                case reflect.Ptr:
                    if !field.IsNil() {
                        isAllNil = false
                    }
                case reflect.Interface:
                    if !field.IsNil() {
                        isAllNil = false
                    }
                }
            }
            return isAllNil
        case reflect.Interface:
            if objElemValue.IsNil() {
                return true
            }
        case reflect.Ptr:
            return self.IsAllNil(objElemValue.Elem())
        default:
            return false
        }
        // objPtrKind := objType.Elem().Kind()
    default:
        switch objKind {
        case reflect.Struct:
            /*
            ** 判断结构体中的所有字段是否都为空
            */
            isAllNil := true
            num := objValue.NumField()
            for i := 0; i < num; i++ {
                field := objValue.Field(i)
                switch field.Type().Kind() {
                case reflect.Ptr:
                    if !field.IsNil() {
                        isAllNil = false
                    }
                case reflect.Interface:
                    if !field.IsNil() {
                        isAllNil = false
                    }
                default:
                }
            }
            return isAllNil
        case reflect.Interface:
            if objValue.IsNil() {
                return true
            }
        default:
        }
    }
    return false
}

func IsAllNil(obj interface{}) bool {
    o := CNilJudge{}
    return o.IsAllNil(obj)
}
