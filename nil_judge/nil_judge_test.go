package nil_judge

import (
    "testing"
    "fmt"
)

type isAllNilJudgeStruct struct {
    F1 *string
    F2 *int
    F3 interface{}
}

func TestIsAllNilJudge(t *testing.T) {
    a := CNilJudge{}
    f1 := "f1"
    s := isAllNilJudgeStruct{
        F1: &f1,
    }
    b := a.IsAllNil(&s)
    fmt.Println(b)
    b = a.IsAllNil(nil)
    fmt.Println(b)
    var i *int
    fmt.Println(a.IsAllNil(&i))
}
