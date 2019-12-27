package strvar_serde

import (
    "testing"
    "fmt"
)

type testStruct struct {
    Name string `field:"name"`
    Age int `field:"age"`
}

func TestCurlyBracesDeserde(t *testing.T) {
    input := "hello, {name}, age: {age}"
    r, err := CurlyBracesDeserde(&input, &testStruct{
        Name: "Jake",
        Age: 20,
    })
    if err != nil {
        return
    }
    fmt.Println(*r)
}

func TestCurlyBracesVarParse(t *testing.T) {
    input := "hello {name}, age: {age}"
    s, err := curlyBracesVarParse(&input, func(v *string)(*string, error) {
        fmt.Println(*v)
        r := "xxx"
        return &r, nil
    })
    if err != nil {
        return
    }
    fmt.Println(*s)
}
