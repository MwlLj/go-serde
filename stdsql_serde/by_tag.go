package stdsql_serde

import (
    "reflect"
    "database/sql"
    "encoding/json"
)

type tagData struct {
    field reflect.Value
    t *string
}

type byTag struct {
}

func (*byTag) assign(rows *sql.Rows, value reflect.Value, t reflect.Type) error {
    num := value.NumField()
    cols := []interface{}{}
    colNames, err := rows.Columns()
    if err != nil {
        return err
    }
    names := map[string]tagData{}
    for i := 0; i < num; i++ {
        field := value.Field(i)
        tag := t.Field(i).Tag
        f := tag.Get(tag_field)
        if f == "" {
            continue
        }
        var ty *string = nil
        typeTag := tag.Get(tag_type)
        if typeTag != "" {
            ty =  &typeTag
        }
        names[f] = tagData{
            field: field,
            t: ty,
        }
    }
    // fmt.Println(colNames)
    cns := []string{}
    for _, colName := range colNames {
        if val, ok := names[colName]; ok {
            field := val.field
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
            if fk == reflect.String ||
                fk == reflect.Struct ||
                fk == reflect.Interface {
                var v sql.NullString
                cols = append(cols, &v)
                cns = append(cns, colName)
            } else if fk == reflect.Int || fk == reflect.Int64 || fk == reflect.Int8 || fk == reflect.Int16 || fk == reflect.Int32 ||
            fk == reflect.Uint8 || fk == reflect.Uint16 || fk == reflect.Uint32 || fk == reflect.Uint64 || fk == reflect.Uint {
                var v sql.NullInt64
                cols = append(cols, &v)
                cns = append(cns, colName)
            } else if fk == reflect.Bool {
                var v sql.NullBool
                cols = append(cols, &v)
                cns = append(cns, colName)
            } else if fk == reflect.Float32 || fk == reflect.Float64 {
                var v sql.NullFloat64
                cols = append(cols, &v)
                cns = append(cns, colName)
            } else {
            }
        }
    }
    if err = rows.Scan(cols...); err != nil {
        return err
    }
    for i, v := range cols {
        n := cns[i]
        if val, ok := names[n]; ok {
            field := val.field
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
                } else if k == reflect.Int || k == reflect.Int64 || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 {
                    va := v.(*sql.NullInt64)
                    if va.Valid {
                        fieldPtrValueElem.SetInt(va.Int64)
                    }
                } else if k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 || k == reflect.Uint {
                    va := v.(*sql.NullInt64)
                    if va.Valid {
                        fieldPtrValueElem.SetUint(uint64(va.Int64))
                    }
                } else if k == reflect.Bool {
                    va := v.(*sql.NullBool)
                    if va.Valid {
                        fieldPtrValueElem.SetBool(va.Bool)
                    }
                } else if k == reflect.Float32 || k == reflect.Float64 {
                    va := v.(*sql.NullFloat64)
                    if va.Valid {
                        fieldPtrValueElem.SetFloat(va.Float64)
                    }
                } else if k == reflect.Struct {
                    /*
                    ** 获取类型, 指定类型序列化
                    */
                    va := v.(*sql.NullString)
                    if va.Valid {
                        if val.t != nil {
                            if *val.t == tag_type_json {
                                fieldValue := reflect.New(fieldPtrType)
                                json.Unmarshal([]byte(va.String), fieldValue.Interface())
                                fieldPtrValueElem.Set(fieldValue.Elem())
                            }
                        } else {
                        }
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
                } else if fk == reflect.Int || fk == reflect.Int64 || fk == reflect.Int8 || fk == reflect.Int16 || fk == reflect.Int32 {
                    field.SetInt(v.(*sql.NullInt64).Int64)
                } else if fk == reflect.Uint8 || fk == reflect.Uint16 || fk == reflect.Uint32 || fk == reflect.Uint64 || fk == reflect.Uint {
                    field.SetUint(uint64(v.(*sql.NullInt64).Int64))
                } else if fk == reflect.Bool {
                    field.SetBool(v.(*sql.NullBool).Bool)
                } else if fk == reflect.Float32 || fk == reflect.Float64 {
                    field.SetFloat(v.(*sql.NullFloat64).Float64)
                } else if fk == reflect.Struct || fk == reflect.Interface {
                    /*
                    ** 获取类型, 指定类型序列化
                    */
                    va := v.(*sql.NullString)
                    if va.Valid {
                        if val.t != nil {
                            if *val.t == tag_type_json {
                                fieldValue := reflect.New(field.Type())
                                json.Unmarshal([]byte(va.String), fieldValue.Interface())
                                field.Set(fieldValue.Elem())
                            }
                        } else {
                        }
                    }
                }
            }
        }
    }
    return nil
}

func ByTag(rows *sql.Rows, output interface{}) error {
    ba := newBase(&byTag{
    })
    return ba.output(rows, output)
}
