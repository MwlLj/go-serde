package strvar_serde

import (
    "fmt"
    "errors"
    "bytes"
    change "github.com/MwlLj/go-serde/obj2map"
)

var _ = fmt.Println

func CurlyBracesDeserde(input *string, obj interface{}) (*string, error) {
    r, err := change.Obj2MapStrStr(obj)
    if err != nil {
        return nil, err
    }
    return curlyBracesDeserde(r, input, nil)
}

func CurlyBracesDeserdeWithCustom(input *string, obj interface{}, customF *func(v *string) (*string, error)) (*string, error) {
    r, err := change.Obj2MapStrStr(obj)
    if err != nil {
        return nil, err
    }
    return curlyBracesDeserde(r, input, customF)
}

/*
** 传入一个字符串, 字符串中的 {var} 使用 结构体 tag 定义的字段的值来替换
*/
func curlyBracesDeserde(r *map[string]string, input *string, customF *func(v *string) (*string, error)) (*string, error) {
    return curlyBracesVarParse(input, func(v *string)(*string, error) {
        var res string
        if v, ok := (*r)[*v]; ok {
            res = v
        } else {
            return nil, errors.New("not found")
        }
        return &res, nil
    }, customF)
}

func CurlyBracesDeserdeMulti(input *string, objs ...interface{}) (*string, error) {
    maps := map[string]string{}
    for _, obj := range objs {
        err := change.Obj2MapStrStrWithCollect(obj, &maps)
        if err != nil {
            return nil, err
        }
    }
    return curlyBracesDeserde(&maps, input, nil)
}

func CurlyBracesDeserdeMultiWithCustom(input *string, customF *func(v *string) (*string, error), objs ...interface{}) (*string, error) {
    maps := map[string]string{}
    for _, obj := range objs {
        err := change.Obj2MapStrStrWithCollect(obj, &maps)
        if err != nil {
            return nil, err
        }
    }
    return curlyBracesDeserde(&maps, input, customF)
}

type mode int8
const (
    _ mode = iota
    normal
    braces
    brackets
)

func curlyBracesVarParse(input *string, f func(v *string) (*string, error), customF *func(v *string) (*string, error)) (*string, error) {
    if input == nil {
        return nil, errors.New("input is nil")
    }
    var result bytes.Buffer 
    var world bytes.Buffer
    m := normal
    for _, c := range *input {
        switch m {
        case normal:
            if c == '{' {
                m = braces
            } else if c == '[' {
                m = brackets
            } else {
                result.WriteRune(c)
            }
        case braces:
            if c == '}' {
                m = normal
                w := world.String()
                r, e := f(&w)
                if e != nil {
                    return nil, e
                }
                result.WriteString(*r)
                world.Reset()
            } else {
                world.WriteRune(c)
            }
        case brackets:
            if c == ']' {
                m = normal
                w := world.String()
                if customF != nil {
                    r, e := (*customF)(&w)
                    if e != nil {
                        return nil, e
                    }
                    result.WriteString(*r)
                }
                world.Reset()
            } else {
                world.WriteRune(c)
            }
        }
    }
    r := result.String()
    return &r, nil
}
