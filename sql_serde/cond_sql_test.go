package sql_serde

import (
    "testing"
    "fmt"
    "bytes"
)

var _ = fmt.Println

type updateSqlObj struct {
    Name *string `field:"name" pos:"0, 1"`
    Likes *string `field:"likes" pos:"0"`
    Age *int `field:"age" pos:"1"`
    Limit *int `field:"limit" pos:"5"`
    Offset *int `field:"offset" pos:"6"`
    Types *[]string `field:"type" pos:"1"`
    BeginTime *string `condfield:"create_time" pos:"2"`
    EndTime *string `field:"endTime" pos:"3" quota:"true"`
    ElementName *string `field:"element_name" pos:"4" quota:"false"`
}

func TestCondSqlSpliceParse(t *testing.T) {
    t.SkipNow()
    sql := "update t_user_info set name = name {,k=v} where name = name {and k = v};"
    s := NewCondSqlSplice()
    s.parse(sql, func(d *data) (string, bool) {
        fmt.Println(d.curIndex, d.kIndex, d.vIndex, d.content)
        return "", true
    })
}

func TestCondSqlAngleParse(t *testing.T) {
    t.SkipNow()
    content := "[prefix] [0-1, 5, 2-4] [ set]"
    s := NewCondSqlSplice()
    var am angleMode = angleModeNormal
    var midClassify string
    var midIndex int = 0
    var word bytes.Buffer
    var ad angleData = angleData{
        prefixIndexGroup: make(map[int]int),
        prefixGroup: make(map[int]prefix),
        splitIndexGroup: make(map[int]int),
        splitGroup: make(map[int]split),
        prefixIndex: 0,
        splitIndex: 0,
    }
    var al angleLast
    for _, c := range content {
        s.angleParse(c, &am, &midClassify, &midIndex, &word, &ad, &al)
    }
    fmt.Println(ad)
}

func TestCondSqlSpliceSerde(t *testing.T) {
    // t.SkipNow()
    // sql := "update t_user_info set name = name {, k = v} where name = name{ and a.k = v}{ and k between v}{ and $v}{ and k like '%v%'}{ limit $v }{ offset $v};"
    // sql := "update t_user_info set name = name {, k = v} where name = name{ and a.k = v} and creatTime > { $v} and creatTime < { $v}{ limit $v }{ offset $v};"
    // sql := "select * from t_vss_vehicle_snapshot_record where 1 = 1{ and k = v}{} and creatTime > { $v} and creatTime < { $v}"
    // sql := "update t_user_info<[prefix] [0] [ set]><[split] [0] [,]>{ k = v}<[prefix] [1-4] [ where]><[split] [1-2, 4] [ and]>{ a.k = v}{ k between v}{ and $v}{ k like '%v%'}{ limit $v}{ offset $v};"
    sql := `
    <[prefix] [0] [ set]>
    <[split] [0] [,]>
    <[prefix] [1-4] [ where]>
    <[split] [1-2, 4] [ and]>
    update t_user_info{ k = v}{ a.k = v}{ k between v}{ and $v}{ k like '%v%'}{ limit $v}{ offset $v};
    `
    // sql := "update t_user_info<[0], set>{<,> k = v};"
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
    elementName := "9.10"
    obj := updateSqlObj{
        Name: &name,
        Age: &age,
        Likes: &likes,
        Limit: &limit,
        Offset: &offset,
        Types: &types,
        BeginTime: &beginTime,
        EndTime: &endTime,
        ElementName: &elementName,
    }
    s := NewCondSqlSplice()
    r, err := s.Serde(sql, &obj)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(*r)
}
