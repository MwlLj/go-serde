package stdsql_serde

import (
    "bytes"
    "strconv"
    "fmt"
    "database/sql"
    "time"
    _ "github.com/go-sql-driver/mysql"
    "testing"
)

type Extra struct {
    F1 string `json:"f1"`
}

type CUserInfo struct {
    Age int `field:"age"`
    Name string `field:"name"`
    Sex *string `field:"sex"`
    Ext Extra `field:"extra" type:"json"`
}

func TestByTag(t *testing.T) {
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
    user := CUserInfo{}
    ba := newBase(&byTag{
    })
    ba.output(rows, &user)
    fmt.Println(user)

    tx.Commit()
}
