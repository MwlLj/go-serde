package obj2map

import (
    "fmt"
    "testing"
)

type obj2MapStrStrTestStruct struct {
    F1 string
    F2 *string `field:"f2"`
    F3 int `field:"f3"`
}

func TestObj2MapStrStr(t *testing.T) {
    f1 := "v1"
    f2 := "v2"
    f3 := 3
    r, err := Obj2MapStrStr(&obj2MapStrStrTestStruct{
        F1: f1,
        F2: &f2,
        F3: f3,
    })
    if err != nil {
        return
    }
    fmt.Println(*r)
}

