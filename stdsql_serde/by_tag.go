package main

import (
    "reflect"
    "errors"
    "bytes"
    "strconv"
    "fmt"
    "database/sql"
    "time"
    _ "github.com/go-sql-driver/mysql"
)

type Tag struct {
    value reflect.Value
    t string
}

func assign(rows *sql.Rows, value reflect.Value, t reflect.Type) error {
    num := value.NumField()
    cols := []interface{}{}
    colNames, err := rows.Columns()
    if err != nil {
        return err
    }
    names := map[string]reflect.Value{}
    for i := 0; i < num; i++ {
        field := value.Field(i)
        f := t.Field(i).Tag.Get("field")
        if f == "" {
            continue
        }
        names[f] = field
    }
    // fmt.Println(colNames)
    cns := []string{}
    for _, colName := range colNames {
        if field, ok := names[colName]; ok {
            fk := field.Kind()
            if fk == reflect.Ptr {
                /*
                ** 指针
                */
                fieldType := field.Type()
                fieldPtrType := fieldType.Elem()
                fk = fieldPtrType.Kind()
            } else {
                /*
                ** 不是指针 => 直接使用fk
                */
            }
            if fk == reflect.String {
                var v sql.NullString
                cols = append(cols, &v)
                cns = append(cns, colName)
            } else if fk == reflect.Int {
                var v sql.NullInt64
                cols = append(cols, &v)
                cns = append(cns, colName)
            }
        }
    }
    rows.Scan(cols...)
    for i, v := range cols {
        n := cns[i]
        if field, ok := names[n]; ok {
            fk := field.Kind()
            if fk == reflect.Ptr {
                fieldType := field.Type()
                fieldValue := reflect.New(fieldType)
                fieldPtrType := fieldType.Elem()
                fieldPtrValue := reflect.New(fieldPtrType)
                fieldPtrValueElem := fieldPtrValue.Elem()
                k := fieldPtrType.Kind()
                if k == reflect.String {
                    va := v.(*sql.NullString)
                    if va.Valid {
                        fieldPtrValueElem.SetString(va.String)
                    }
                } else if k == reflect.Int {
                    va := v.(*sql.NullInt64)
                    if va.Valid {
                        fieldPtrValueElem.SetInt(va.Int64)
                    }
                }
                fieldValue.Elem().Set(fieldPtrValue)
                field.Set(fieldValue.Elem())
            } else {
                /*
                ** 不是指针
                */
                if fk == reflect.String {
                    field.SetString(v.(*sql.NullString).String)
                } else if fk == reflect.Int {
                    field.SetInt(v.(*sql.NullInt64).Int64)
                }
            }
        }
    }
    /*
    colLen := len(cols)
    for i := 0; i < num; i++ {
        field := value.Field(i)
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
    */
    return nil
}

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
        // fmt.Println("is slice")
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
                    assign(rows, slicePtrValueElem, slicePtrType)
                    /*
                    num := slicePtrValueElem.NumField()
                    cols := []interface{}{}
                    colNames, err := rows.Columns()
                    if err != nil {
                        return err
                    }
                    names := map[string]reflect.Value{}
                    for i := 0; i < num; i++ {
                        field := slicePtrValueElem.Field(i)
                        f := slicePtrType.Field(i).Tag.Get("field")
                        if f == "" {
                            continue
                        }
                        names[f] = field
                    }
                    for _, colName := range colNames {
                        if field, ok := names[colName]; ok {
                            fk := field.Kind()
                            if fk == reflect.String {
                                var v sql.NullString
                                cols = append(cols, &v)
                            } else if fk == reflect.Int {
                                var v sql.NullInt64
                                cols = append(cols, &v)
                            }
                        } else {
                            cols = append(cols, &sql.NullString{})
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
                    assign(rows, slicePtrValueElem, slicePtrType)
                    /*
                    fk := slicePtrValueElem.Kind()
                    if fk == reflect.String {
                        slicePtrValueElem.SetString("Mike")
                    } else if fk == reflect.Int {
                        slicePtrValueElem.SetInt(10)
                    }
                    */
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
                    assign(rows, sliceValueElem, sliceType)
                    /*
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
                    */
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                } else {
                    // fmt.Println("sliceType is not struct")
                    /*
                    ** 读取一列
                    ** []string{}
                    */
                    sliceValue := reflect.New(sliceType)
                    sliceValueElem := sliceValue.Elem()
                    assign(rows, sliceValueElem, sliceType)
                    /*
                    fk := sliceValueElem.Kind()
                    if fk == reflect.String {
                        sliceValueElem.SetString("Alis")
                    } else if fk == reflect.Int {
                        sliceValueElem.SetInt(10)
                    }
                    */
                    outValue.Set(reflect.Append(outValue, sliceValueElem))
                }
            }
        }
    } else {
        // fmt.Println("is not slice")
        objType := outValue.Type()
        /*
        ** 判断类型是否是结构体
        */
        if objType.Kind() == reflect.Struct {
            // fmt.Println("objType is struct")
            /*
            ** 读取多列
            */
            sliceValue := reflect.New(objType)
            sliceValueElem := sliceValue.Elem()
            assign(rows, sliceValueElem, objType)
            /*
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
            */
            outValue.Set(sliceValueElem)
        } else {
            // fmt.Println("objType is not struct")
            /*
            ** 读取一列
            */
            sliceValue := reflect.New(objType)
            sliceValueElem := sliceValue.Elem()
            assign(rows, sliceValueElem, objType)
            /*
            fmt.Println(sliceValueElem.Type())
            fk := sliceValueElem.Kind()
            if fk == reflect.String {
                sliceValueElem.SetString("Blue")
            } else if fk == reflect.Int {
                sliceValueElem.SetInt(10)
            }
            */
            outValue.Set(sliceValueElem)
        }
        /*
        ** 只读取一行
        */
    }
    return nil
}

type CUserInfo struct {
    Age int `field:"age"`
    Name string `field:"name"`
    Sex *string `field:"sex"`
}

func main() {
    b := bytes.Buffer{}
    b.WriteString("root")
    b.WriteString(":")
    b.WriteString("123456")
    b.WriteString("@tcp(")
    b.WriteString("127.0.0.1")
    b.WriteString(":")
    b.WriteString(strconv.FormatUint(uint64(3306), 10))
    b.WriteString(")/")
    b.WriteString("test")
    db, err := sql.Open("mysql", b.String())
    if err != nil {
        fmt.Println(err)
        return
    }
    defer db.Close()
    db.SetMaxOpenConns(2000)
    db.SetMaxIdleConns(1000)
    db.SetConnMaxLifetime(time.Second * 10)
    db.Ping()
    tx, err := db.Begin()
    if err != nil {
        return
    }
    rows, err := db.Query(fmt.Sprintf(`select * from t_user_info;`))
    if err != nil {
        tx.Rollback()
        return
    }
    defer rows.Close()

    user := []*CUserInfo{}
    output(rows, &user)
    for _, u := range user {
        if u.Sex != nil {
            fmt.Println(u.Age, u.Name, *u.Sex)
        } else {
            fmt.Println(u.Age, u.Name)
        }
    }
    /*
    user1 := CUserInfo{}
    fmt.Println("------CUserInfo{}------")
    output(rows, &user1)
    fmt.Println(user1)
    user2 := []CUserInfo{}
    fmt.Println("------[]CUserInfo{}------")
    output(rows, &user2)
    for _, user := range user2 {
        fmt.Println(user)
    }
    user3 := []*CUserInfo{}
    fmt.Println("------[]*CUserInfo{}------")
    output(rows, &user3)
    for _, user := range user3 {
        fmt.Println(user)
    }
    user4 := []string{}
    fmt.Println("------[]string{}------")
    output(rows, &user4)
    for _, user := range user4 {
        fmt.Println(user)
    }
    user5 := []*string{}
    fmt.Println("------[]*string{}------")
    output(rows, &user5)
    for _, user := range user5 {
        fmt.Println(*user)
    }
    var user6 string
    fmt.Println("------string{}------")
    output(rows, &user6)
    fmt.Println(user6)
    */

    tx.Commit()
}
