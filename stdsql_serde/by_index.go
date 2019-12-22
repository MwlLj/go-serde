package stdsql_serde

import (
    "reflect"
    "errors"
    "database/sql"
    "encoding/json"
)

type byIndex struct {
}

func (*byIndex) assign(rows *sql.Rows, value reflect.Value, t reflect.Type) error {
    num := value.NumField()
    cols := []interface{}{}
    for i := 0; i < num; i++ {
        field := value.Field(i)
        fk := field.Kind()
        if fk == reflect.Ptr {
            fieldType := field.Type()
            fieldPtrType := fieldType.Elem()
            fk = fieldPtrType.Kind()
        } else {
        }
        if fk == reflect.String || fk == reflect.Struct {
            var v sql.NullString
            cols = append(cols, &v)
        } else if fk == reflect.Int || fk == reflect.Int64 || fk == reflect.Int8 || fk == reflect.Int16 || fk == reflect.Int32 ||
        fk == reflect.Uint8 || fk == reflect.Uint16 || fk == reflect.Uint32 || fk == reflect.Uint64 || fk == reflect.Uint {
            var v sql.NullInt64
            cols = append(cols, &v)
        } else if fk == reflect.Bool {
            var v sql.NullBool
            cols = append(cols, &v)
        } else if fk == reflect.Float32 || fk == reflect.Float64 {
            var v sql.NullFloat64
            cols = append(cols, &v)
        } else {
        }
    }
    if err := rows.Scan(cols...); err != nil {
        return err
    }
    colLen := len(cols)
    if colLen != num {
        return errors.New("cols not match")
    }
    for i := 0; i < num; i++ {
        if i + 1 > colLen {
            return errors.New("cols not match, check tags")
        }
        field := value.Field(i)
        fk := field.Kind()
        switch fk {
        case reflect.Ptr: {
            fieldType := field.Type()
            fieldPtrType := fieldType.Elem()
            fieldPtrValue := reflect.New(fieldPtrType)
            fieldPtrValueElem := fieldPtrValue.Elem()
            k := fieldPtrType.Kind()
            v := cols[i]
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
                tag := t.Field(i).Tag
                typeTag := tag.Get(tag_type)
                va := v.(*sql.NullString)
                if va.Valid {
                    if typeTag != "" {
                        if typeTag == tag_type_json {
                            fieldValue := reflect.New(fieldPtrType)
                            json.Unmarshal([]byte(va.String), fieldValue.Interface())
                            fieldPtrValueElem.Set(fieldValue.Elem())
                        }
                    } else {
                    }
                }
            }
            field.Set(fieldPtrValue)
        }
        default: {
            v := cols[i]
            if fk == reflect.String {
                field.SetString(v.(*sql.NullString).String)
            } else if fk == reflect.Int || fk == reflect.Int64 || fk == reflect.Int8 || fk == reflect.Int16 || fk == reflect.Int32 {
                field.SetInt(v.(*sql.NullInt64).Int64)
            } else if fk == reflect.Uint8 || fk == reflect.Uint16 || fk == reflect.Uint32 || fk == reflect.Uint64 || fk == reflect.Uint {
                field.SetUint(uint64(v.(sql.NullInt64).Int64))
            } else if fk == reflect.Bool {
                field.SetBool(v.(*sql.NullBool).Bool)
            } else if fk == reflect.Float32 || fk == reflect.Float64 {
                field.SetFloat(v.(*sql.NullFloat64).Float64)
            } else if fk == reflect.Struct {
                tag := t.Field(i).Tag
                typeTag := tag.Get(tag_type)
                va := v.(*sql.NullString)
                if va.Valid {
                    if typeTag != "" {
                        if typeTag == tag_type_json {
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

func ByIndex(rows *sql.Rows, output interface{}) error {
    ba := newBase(&byIndex{
    })
    return ba.output(rows, output)
}
