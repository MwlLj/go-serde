package stdsql_serde

import (
    "reflect"
    "errors"
    "fmt"
    "database/sql"
    // "unsafe"
)

var _ = fmt.Println

const (
    tag_field string = "field"
    tag_type string = "type"
    tag_type_json string = "json"
)

type ISet interface {
    assign(rows *sql.Rows, value reflect.Value, t reflect.Type) error
}

type base struct {
    set ISet
}

func (self *base) output(rows *sql.Rows, out interface{}) (bool, error) {
    if out == nil {
        return false, nil
    }
    outValuePtr := reflect.ValueOf(out)
    if outValuePtr.IsNil() {
        return false, nil
    }
    /*
    ** 判断是否可以设置
    */
    if outValuePtr.Kind() != reflect.Ptr {
        return false, errors.New("can not set, please use pointer")
    }
    /*
    ** 取出指针的值
    */
    outValue := outValuePtr.Elem()
    /*
    ** 判断是否是 slice
    */
    if outValue.Kind() == reflect.Slice {
        // fmt.Println("is slice")
        /*
        ** 读取每一个行的值, 然后追加到 slice 中
        */
        has := false
        for rows.Next() {
            has = true
            /*
            ** 获取 slice 中的类型
            */
            sliceType := outValue.Type().Elem()
            /*
            ** 判断sliceType是否是指针
            */
            if sliceType.Kind() == reflect.Ptr {
                // fmt.Println("slice type is ptr")
                /*
                ** 获取sliceType指针类型中的实体类型
                */
                sliceValue := reflect.New(sliceType)
                sliceValueElem := sliceValue.Elem()
                slicePtrType := sliceType.Elem()
                // fmt.Printf("slicePtrType: %v\n", slicePtrType.String())
                /*
                ** 判断实体类型是否是结构体
                */
                if slicePtrType.Kind() == reflect.Struct {
                    // fmt.Println("slicePtrType is struct")
                    /*
                    ** 读取多列
                    ** []*CUserInfo{}
                    */
                    slicePtrValue := reflect.New(slicePtrType)
                    slicePtrValueElem := slicePtrValue.Elem()
                    self.set.assign(rows, slicePtrValueElem, slicePtrType)
                    sliceValueElem.Set(slicePtrValue)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                } else {
                    // fmt.Println("slicePtrType is not struct")
                    /*
                    ** 读取一列
                    ** []*string
                    */
                    slicePtrValue := reflect.New(slicePtrType)
                    slicePtrValueElem := slicePtrValue.Elem()
                    self.set.assign(rows, slicePtrValueElem, slicePtrType)
                    sliceValueElem.Set(slicePtrValue)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                }
            } else {
                // fmt.Println("slice type is not ptr")
                /*
                ** 判断类型是否是结构体
                */
                if sliceType.Kind() == reflect.Struct {
                    // fmt.Println("sliceType is struct")
                    /*
                    ** 读取多列
                    ** []CUserInfo{}
                    */
                    sliceValue := reflect.New(sliceType)
                    sliceValueElem := sliceValue.Elem()
                    self.set.assign(rows, sliceValueElem, sliceType)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                } else {
                    // fmt.Println("sliceType is not struct")
                    /*
                    ** 读取一列
                    ** []string{}
                    */
                    sliceValue := reflect.New(sliceType)
                    sliceValueElem := sliceValue.Elem()
                    self.set.assign(rows, sliceValueElem, sliceType)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                }
            }
        }
        if !has {
            return false, nil
        }
    } else {
        // fmt.Println("is not slice")
        /*
        ** 只读取一行
        */
        objType := outValue.Type()
        switch objType.Kind() {
        case reflect.Ptr:
            sliceValue := reflect.New(objType.Elem())
            isExist, err := self.output(rows, sliceValue.Interface())
            if err != nil {
                return false, err
            }
            if isExist {
                outValue.Set(sliceValue)
            }
        default:
            if rows.Next() {
                // fmt.Println("is not slice")
                /*
                ** 判断类型是否是结构体
                */
                switch objType.Kind() {
                case reflect.Struct:
                    // fmt.Println("objType is struct")
                    /*
                    ** 读取多列
                    */
                    sliceValue := reflect.New(objType)
                    sliceValueElem := sliceValue.Elem()
                    self.set.assign(rows, sliceValueElem, objType)
                    outValue.Set(sliceValueElem)
                default:
                    // fmt.Println("objType is not struct")
                    /*
                    ** 读取一列
                    */
                    sliceValue := reflect.New(objType)
                    sliceValueElem := sliceValue.Elem()
                    self.set.assign(rows, sliceValueElem, objType)
                    outValue.Set(sliceValueElem)
                }
            } else {
                return false, nil
            }
        }
    }
    return true, nil
}

func newBase(set ISet) *base {
    return &base{
        set: set,
    }
}
