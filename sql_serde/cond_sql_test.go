package sql_serde

import (
    "testing"
    "fmt"
)

var _ = fmt.Println

type updateSqlObj struct {
    Name *string `field:"name" pos:"0"`
    Likes *string `field:"likes" pos:"0"`
    Age *int `field:"age" pos:"1"`
}

func TestCondSqlSpliceParse(t *testing.T) {
    sql := "update t_user_info set name = name {,k=v} where name = name {and k = v};"
    s := NewCondSqlSplice()
    s.parse(sql, func(d *data) string {
        fmt.Println(d.curIndex, d.kIndex, d.vIndex, d.content)
        return ""
    })
}

func TestCondSqlSpliceSerde(t *testing.T) {
    sql := "update t_user_info set name = name {, k = v} where name = name {and k = v};"
    name := "jake"
    age := 20
    likes := "fruit"
    obj := updateSqlObj{
        Name: &name,
        Age: &age,
        Likes: &likes,
    }
    s := NewCondSqlSplice()
    r, err := s.Serde(sql, &obj)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(*r)
}
