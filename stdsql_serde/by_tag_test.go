package stdsql_serde

import (
    "bytes"
    "strconv"
    "fmt"
    "database/sql"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "testing"
    // "reflect"
)

type Extra struct {
    F1 string `json:"f1"`
}

type CUserInfo struct {
    Age int `field:"age"`
    Name string `field:"name"`
    Sex *string `field:"sex"`
    Ext *Extra `field:"extra"`
}

type CUserInfo2 struct {
    Age int `field:"age"`
    Name string `field:"name"`
    Sex *string `field:"sex"`
    Ext interface{} `field:"extra"`
}

func TestByTag(t *testing.T) {
    t.SkipNow()
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
    // rows, err := db.Query(fmt.Sprintf(`select * from t_user_info;`))
    rows, err := db.Query(fmt.Sprintf(`select count(0) from t_user_info;`))
    if err != nil {
        tx.Rollback()
        return
    }
    defer rows.Close()

    /*
    user := []*CUserInfo{}
    output(rows, &user)
    for _, u := range user {
        if u.Sex != nil {
            fmt.Println(u.Age, u.Name, *u.Sex, u.Ext)
        } else {
            fmt.Println(u.Age, u.Name)
        }
    }
    */
    // user := CUserInfo{}
    // ByTag(rows, &user)
    // fmt.Println(user.Name, user.Age, *user.Sex, *user.Ext)

    /*
    var user *CUserInfo
    err = ByTag(rows, &user)
    if user != nil {
        fmt.Println(user.Name, user.Age, user.Sex, user.Ext)
    } else {
        fmt.Println("not found")
    }
    */

    var count *int
    err = ByTag(rows, &count)
    fmt.Println(*count)

    tx.Commit()
}

func TestByTagFromTemplateObj(t *testing.T) {
    t.SkipNow()
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

    /*
    user := []*CUserInfo{}
    output(rows, &user)
    for _, u := range user {
        if u.Sex != nil {
            fmt.Println(u.Age, u.Name, *u.Sex, u.Ext)
        } else {
            fmt.Println(u.Age, u.Name)
        }
    }
    */
    ext := Extra{}
    // typ := reflect.TypeOf(ext)
    // n := reflect.New(typ)
    user := CUserInfo2{
        Ext: ext,
    }
    ByTag(rows, &user)
    fmt.Println(user.Name, user.Age, user.Sex, user.Ext)

    tx.Commit()
}

func TestByTagWithValues(t *testing.T) {
    t.SkipNow()
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

    /*
    user := []*CUserInfo{}
    output(rows, &user)
    for _, u := range user {
        if u.Sex != nil {
            fmt.Println(u.Age, u.Name, *u.Sex, u.Ext)
        } else {
            fmt.Println(u.Age, u.Name)
        }
    }
    */
    // typ := reflect.TypeOf(ext)
    // n := reflect.New(typ)
    users := []CUserInfo2{
    }
    values := map[string]interface{}{}
    ext := Extra{}
    // v := reflect.ValueOf(&ext).Elem()
    // fmt.Println(v.CanAddr())
    values["extra"] = &ext
    ByTagWithValues(rows, &users, values)
    for _, user := range users {
    	fmt.Println(user)
    }

    tx.Commit()
}

type PersonSnapshotResponse struct {
    ChannelName *string `json:"channelName" field:"channelName" pos:"0"`
    PanoramaCapturePhotoUrl *string `json:"panoramaCapturePhotoUrl" field:"panoramaCapturePhotoUrl" pos:"0"`
    FaceCapturePhotoUrl *string `json:"faceCapturePhotoUrl" field:"faceCapturePhotoUrl" pos:"0"`
    FaceSimilarity *int `json:"faceSimilarity" field:"faceSimilarity" pos:"0"`
    CaptureDateTime *string `json:"captureDateTime" field:"captureDateTime" pos:"0"`
}

func TestByTag2(t *testing.T) {
    // t.SkipNow()
    b := bytes.Buffer{}
    b.WriteString("root")
    b.WriteString(":")
    b.WriteString("nYd*4M]ipIv+")
    b.WriteString("@tcp(")
    b.WriteString("192.168.9.238")
    b.WriteString(":")
    b.WriteString(strconv.FormatUint(uint64(3306), 10))
    b.WriteString(")/")
    b.WriteString("iacdb_4f186c0f38d8429898b2dcac72de43c0")
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
    fmt.Println("after connect")
    tx, err := db.Begin()
    if err != nil {
        return
    }
    // rows, err := db.Query(fmt.Sprintf(`select * from t_user_info;`))
    rows, err := db.Query(fmt.Sprintf(`
        select * from (select b.channelName,b.panoramaCapturePhotoUrl,b.faceCapturePhotoUrl,b.faceSimilarity,b.captureDateTime from t_vss_face_info a,(select faceUuid,channelName,panoramaCapturePhotoUrl,faceCapturePhotoUrl,faceSimilarity,faceCapturePhotoQuality,captureDateTime from t_vss_facerecognition_record_2020_02 union select faceUuid,channelName,panoramaCapturePhotoUrl,faceCapturePhotoUrl,faceSimilarity,faceCapturePhotoQuality,captureDateTime from t_vss_facerecognition_record_2020_03) b  where a.faceUuid = b.faceUuid  and a.staffUuid = 'a1bb1e7ea84043c49e5194324326d3af' and b.faceCapturePhotoQuality = 'HIGH' and b.captureDateTime >= '2020-01-03 00:00:00' and b.captureDateTime <= '2020-03-03 18:52:56') _paging limit 0,5
        `))
    if err != nil {
        tx.Rollback()
        return
    }
    defer rows.Close()

    /*
    user := []*CUserInfo{}
    output(rows, &user)
    for _, u := range user {
        if u.Sex != nil {
            fmt.Println(u.Age, u.Name, *u.Sex, u.Ext)
        } else {
            fmt.Println(u.Age, u.Name)
        }
    }
    */
    // user := CUserInfo{}
    // ByTag(rows, &user)
    // fmt.Println(user.Name, user.Age, *user.Sex, *user.Ext)

    /*
    var user *CUserInfo
    err = ByTag(rows, &user)
    if user != nil {
        fmt.Println(user.Name, user.Age, user.Sex, user.Ext)
    } else {
        fmt.Println("not found")
    }
    */

    var result []PersonSnapshotResponse
    err = ByTag(rows, &result)
    for _, v := range result {
        fmt.Println(*v.ChannelName, *v.PanoramaCapturePhotoUrl, *v.FaceCapturePhotoUrl, *v.FaceSimilarity, *v.CaptureDateTime)
    }

    tx.Commit()
}
