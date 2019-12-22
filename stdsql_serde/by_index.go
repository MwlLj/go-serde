package stdsql_serde

import (
    "reflect"
    "errors"
    "fmt"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func output(rows *sql.Rows, out interface{}) error {
    outValuePtr := reflect.ValueOf(out)
    /*
    ** 判断是否可以设置
    */
    if outValuePtr.Kind() != reflect.Ptr {
        return errors.New("can not set, please use pointer")
    }
    /*
    ** 取出指针的值
    */
    outValue := outValuePtr.Elem()
    /*
    ** 判断是否是 slice
    */
    if outValue.Kind() == reflect.Slice {
        fmt.Println("is slice")
        /*
        ** 读取每一个行的值, 然后追加到 slice 中
        */
        for rows.Next() {
            /*
            ** 获取 slice 中的类型
            */
            sliceType := outValue.Type().Elem()
            /*
            ** 判断sliceType是否是指针
            */
            if sliceType.Kind() == reflect.Ptr {
                fmt.Println("slice type is ptr")
                /*
                ** 获取sliceType指针类型中的实体类型
                */
                sliceValue := reflect.New(sliceType)
                sliceValueElem := sliceValue.Elem()
                slicePtrType := sliceType.Elem()
                fmt.Printf("slicePtrType: %v\n", slicePtrType.String())
                /*
                ** 判断实体类型是否是结构体
                */
                if slicePtrType.Kind() == reflect.Struct {
                    fmt.Println("slicePtrType is struct")
                    /*
                    ** 读取多列
                    ** []*CUserInfo{}
                    */
                    slicePtrValue := reflect.New(slicePtrType)
                    slicePtrValueElem := slicePtrValue.Elem()
                    num := slicePtrValueElem.NumField()
                    cols := []interface{}{}
                    for i := 0; i < num; i++ {
                        field := slicePtrValueElem.Field(i)
                        fk := field.Kind()
                        if fk == reflect.String {
                            var v sql.NullString
                            cols = append(cols, &v)
                        } else if fk == reflect.Int {
                            var v sql.NullInt64
                            cols = append(cols, &v)
                        }
                    }
                    rows.Scan(cols...)
                    colLen := len(cols)
                    for i := 0; i < num; i++ {
                        field := slicePtrValueElem.Field(i)
                        fk := field.Kind()
                        if i + 1 > colLen {
                            break
                        }
                        v := cols[i]
                        if fk == reflect.String {
                            field.SetString(v.(*sql.NullString).String)
                        } else if fk == reflect.Int {
                            field.SetInt(v.(*sql.NullInt64).Int64)
                        }
                    }
                    /*
                    for _, col := range cols {
                        fmt.Printf("-----%v----- ", col)
                    }
                    */
                    /*
                    colValues, err := rows.Columns()
                    if err != nil {
                        return err
                    }
                    colTypes, err := rows.ColumnTypes()
                    if err != nil {
                        return nil
                    }
                    var _ = colTypes
                    colLen := len(colValues)
                    for i := 0; i < num; i++ {
                        field := slicePtrValueElem.Field(i)
                        fk := field.Kind()
                        if i + 1 > colLen {
                            break
                        }
                        v := colValues[i]
                        if fk == reflect.String {
                            field.SetString(v)
                            fmt.Println(v)
                        } else if fk == reflect.Int {
                            fmt.Println(v)
                            iv, err := strconv.ParseInt(v, 10, 64)
                            if err != nil {
                                return errors.New(fmt.Sprintf("field: %s is not int", field.Type().Name))
                            }
                            field.SetInt(iv)
                        }
                    }
                    */
                    sliceValueElem.Set(slicePtrValue)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                } else {
                    fmt.Println("slicePtrType is not struct")
                    /*
                    ** 读取一列
                    ** []*string
                    */
                    slicePtrValue := reflect.New(slicePtrType)
                    slicePtrValueElem := slicePtrValue.Elem()
                    fk := slicePtrValueElem.Kind()
                    if fk == reflect.String {
                        slicePtrValueElem.SetString("Mike")
                    } else if fk == reflect.Int {
                        slicePtrValueElem.SetInt(10)
                    }
                    sliceValueElem.Set(slicePtrValue)
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                }
            } else {
                fmt.Println("slice type is not ptr")
                /*
                ** 判断类型是否是结构体
                */
                if sliceType.Kind() == reflect.Struct {
                    fmt.Println("sliceType is struct")
                    /*
                    ** 读取多列
                    ** []CUserInfo{}
                    */
                    sliceValue := reflect.New(sliceType)
                    sliceValueElem := sliceValue.Elem()
                    num := sliceValueElem.NumField()
                    for i := 0; i < num; i++ {
                        field := sliceValueElem.Field(i)
                        fk := field.Kind()
                        if fk == reflect.String {
                            field.SetString("Lan")
                        } else if fk == reflect.Int {
                            field.SetInt(21)
                        }
                    }
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                } else {
                    fmt.Println("sliceType is not struct")
                    /*
                    ** 读取一列
                    ** []string{}
                    */
                    sliceValue := reflect.New(sliceType)
                    sliceValueElem := sliceValue.Elem()
                    fk := sliceValueElem.Kind()
                    if fk == reflect.String {
                        sliceValueElem.SetString("Alis")
                    } else if fk == reflect.Int {
                        sliceValueElem.SetInt(10)
                    }
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                }
            }
        }
    } else {
        fmt.Println("is not slice")
        objType := outValue.Type()
        /*
        ** 判断类型是否是结构体
        */
        if objType.Kind() == reflect.Struct {
            fmt.Println("objType is struct")
            /*
            ** 读取多列
            */
            sliceValue := reflect.New(objType)
            sliceValueElem := sliceValue.Elem()
            num := sliceValueElem.NumField()
            for i := 0; i < num; i++ {
                field := sliceValueElem.Field(i)
                fk := field.Kind()
                if fk == reflect.String {
                    field.SetString("Red")
                } else if fk == reflect.Int {
                    field.SetInt(21)
                }
            }
            outValue.Set(sliceValueElem)
        } else {
            fmt.Println("objType is not struct")
            /*
            ** 读取一列
            */
            sliceValue := reflect.New(objType)
            sliceValueElem := sliceValue.Elem()
            fmt.Println(sliceValueElem.Type())
            fk := sliceValueElem.Kind()
            if fk == reflect.String {
                sliceValueElem.SetString("Blue")
            } else if fk == reflect.Int {
                sliceValueElem.SetInt(10)
            }
            outValue.Set(sliceValueElem)
        }
        /*
        ** 只读取一行
        */
    }
    return nil
}
