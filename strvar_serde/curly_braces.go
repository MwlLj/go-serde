package strvar_serde

import (
    "fmt"
    "errors"
    "bytes"
    change "dbserver/parse/serde/obj2map"
)

var _ = fmt.Println

func CurlyBracesDeserde(input *string, obj interface{}) (*string, error) {
    r, err := change.Obj2MapStrStr(obj)
    if err != nil {
        return nil, err
    }
    return curlyBracesDeserde(r, input, nil, nil)
}

func CurlyBracesDeserdeWithCustom(input *string, obj interface{}, customF *func(v *string) (*string, error)) (*string, error) {
    r, err := change.Obj2MapStrStr(obj)
    if err != nil {
        return nil, err
    }
    return curlyBracesDeserde(r, input, nil, customF)
}

/*
** 传入一个字符串, 字符串中的 {var} 使用 结构体 tag 定义的字段的值来替换
*/
func curlyBracesDeserde(r *map[string]string, input *string, areaF func(k *string, v *string), customF *func(v *string) (*string, error)) (*string, error) {
    return curlyBracesVarParse(input, func(v *string)(*string, error) {
        var res string
        if v, ok := (*r)[*v]; ok {
            res = v
        } else {
            return nil, errors.New("not found")
        }
        return &res, nil
    }, areaF, customF)
}

func CurlyBracesDeserdeMulti(input *string, objs ...interface{}) (*string, error) {
    maps := map[string]string{}
    for _, obj := range objs {
        err := change.Obj2MapStrStrWithCollect(obj, &maps)
        if err != nil {
            return nil, err
        }
    }
    return curlyBracesDeserde(&maps, input, nil, nil)
}

func CurlyBracesDeserdeMultiWithCustom(input *string, customF *func(v *string) (*string, error), objs ...interface{}) (*string, error) {
    maps := map[string]string{}
    for _, obj := range objs {
        err := change.Obj2MapStrStrWithCollect(obj, &maps)
        if err != nil {
            return nil, err
        }
    }
    return curlyBracesDeserde(&maps, input, nil, customF)
}

func CurlyBracesDeserdeMultiWithCustomAndArea(input *string, areaF func(k *string, v *string), customF *func(v *string) (*string, error), objs ...interface{}) (*string, error) {
    maps := map[string]string{}
    for _, obj := range objs {
        err := change.Obj2MapStrStrWithCollect(obj, &maps)
        if err != nil {
            return nil, err
        }
    }
    return curlyBracesDeserde(&maps, input, areaF, customF)
}

type mode int8
const (
    _ mode = iota
    normal
    braces
    brackets
    angle_brackets
    angle_brackets_key
    angle_brackets_value
)

func curlyBracesVarParse(input *string, f func(v *string) (*string, error), areaF func(k *string, v *string), customF *func(v *string) (*string, error)) (*string, error) {
    if input == nil {
        return nil, errors.New("input is nil")
    }
    var result bytes.Buffer 
    var world bytes.Buffer
    var angleBracketsKey bytes.Buffer
    m := normal
    angleBracketsMode := angle_brackets_key
    length := len(*input)
    loop:
    for i, c := range *input {
        switch m {
        case normal:
            if c == '{' {
                m = braces
            } else if c == '[' {
                m = brackets
            } else if c == '<' {
                m = angle_brackets;
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
        case angle_brackets:
            switch angleBracketsMode {
            case angle_brackets_key:
                if c == '>' {
                    if i+1 > length - 1 {
                        break loop
                    }
                    if (*input)[i+1] == '<' {
                        angleBracketsMode = angle_brackets_value
                    } else {
                        break loop
                    }
                } else {
                    angleBracketsKey.WriteRune(c)
                }
            case angle_brackets_value:
                if c == '>' {
                    m = normal
                    angleBracketsMode = angle_brackets_key
                    in := world.String()
                    v, e := curlyBracesVarParse(&in, f, areaF, customF)
                    if e != nil {
                        break loop
                    }
                    result.WriteString(*v)
                    if areaF != nil {
                        k := angleBracketsKey.String()
                        areaF(&k, v)
                    }
                    world.Reset()
                    angleBracketsKey.Reset()
                } else {
                    if c == '<' {
                        continue
                    }
                    world.WriteRune(c)
                }
            }
        }
    }
    r := result.String()
    return &r, nil
}
