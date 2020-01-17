package sql_serde

import (
    "testing"
    "fmt"
)

var _ = fmt.Println

type updateSqlObj struct {
    Name *string `field:"name" pos:"0, 1"`
    Likes *string `field:"likes" pos:"0"`
    Age *int `field:"age" pos:"1"`
    Limit *int `field:"limit" pos:"4"`
    Offset *int `field:"offset" pos:"5"`
    Types *[]string `field:"type" pos:"1"`
    BeginTime *string `condfield:"create_time" pos:"2"`
    EndTime *string `field:"endTime" pos:"3" quota:"true"`
}

func TestCondSqlSpliceParse(t *testing.T) {
    t.SkipNow()
    sql := "update t_user_info set name = name {,k=v} where name = name {and k = v};"
    s := NewCondSqlSplice()
    s.parse(sql, func(d *data) string {
        fmt.Println(d.curIndex, d.kIndex, d.vIndex, d.content)
        return ""
    })
}

func TestCondSqlSpliceSerde(t *testing.T) {
    sql := "update t_user_info set name = name {, k = v} where name = name{ and a.k = v}{ k between v}{ and $v}{ limit $v }{ offset $v};"
    // sql := "update t_user_info set name = name {, k = v} where name = name{ and a.k = v} and creatTime > { $v} and creatTime < { $v}{ limit $v }{ offset $v};"
    // sql := "select * from t_vss_vehicle_snapshot_record where 1 = 1{ and k = v}{} and creatTime > { $v} and creatTime < { $v}"
    name := "jake"
    age := 20
    likes := "fruit"
    limit := 1
    offset := 0
    types := []string{
        "v1", "v2",
    }
    beginTime := "2020-01-09 00:00:00"
    endTime := "2020-01-09 23:59:59"
    obj := updateSqlObj{
        Name: &name,
        Age: &age,
        Likes: &likes,
        Limit: &limit,
        Offset: &offset,
        Types: &types,
        BeginTime: &beginTime,
        EndTime: &endTime,
    }
    s := NewCondSqlSplice()
    r, err := s.Serde(sql, &obj)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(*r)
}
