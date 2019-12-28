package strvar_serde

import (
    "testing"
    "fmt"
    "errors"
)

type testStruct struct {
    Name string `field:"name"`
    Age int `field:"age"`
}

type testStruct1 struct {
    Name string `field:"name"`
}

type testStruct2 struct {
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

func TestCurlyBracesDeserdeWithCustom(t *testing.T) {
    input := "hello, {name}, age: {age}, [uuid]"
    customF := func(v *string) (*string, error) {
        if *v == "uuid" {
            uuid := "123456"
            return &uuid, nil
        }
        return nil, errors.New("not match")
    }
    r, err := CurlyBracesDeserdeWithCustom(&input, &testStruct{
        Name: "Jake",
        Age: 20,
    }, &customF)
    if err != nil {
        return
    }
    fmt.Println(*r)
}

func TestCurlyBracesDeserdeMulti(t *testing.T) {
    input := "hello, {name}, age: {age}"
    r, err := CurlyBracesDeserdeMulti(&input, &testStruct1{
        Name: "Jake",
    }, &testStruct2{
        Age: 20,
    })
    if err != nil {
        return
    }
    fmt.Println(*r)
}

func TestCurlyBracesDeserdeMultiWithCustom(t *testing.T) {
    input := "hello, {name}, age: {age}, [uuid]"
    customF := func(v *string) (*string, error) {
        if *v == "uuid" {
            uuid := "123456"
            return &uuid, nil
        }
        return nil, errors.New("not match")
    }
    r, err := CurlyBracesDeserdeMultiWithCustom(&input, &customF, &testStruct1{
        Name: "Jake",
    }, &testStruct2{
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
    }, nil)
    if err != nil {
        return
    }
    fmt.Println(*s)
}
