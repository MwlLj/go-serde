package sql_serde

import (
    "reflect"
    "errors"
    "strconv"
    "bytes"
    "fmt"
)

var _ = fmt.Println

const (
    cond_tag_field string = "field"
    cond_tag_cond_field string = "condfield"
    cond_tag_pos string = "pos"
    cond_tag_quota string = "quota"
    keyword_key string = "k"
    keyword_value string = "v"
    keyword_va string = "$v"
    keyword_quota_true string = "true"
    keyword_quota_false string = "false"
)

var (
    keyword_key_len int = len(keyword_key)
    keyword_value_len int = len(keyword_value)
    keyword_va_len int = len(keyword_va)
)

type Mode int8
const (
    _ Mode = iota
    normal
    bigBrackets
)

type blockMode int8
const (
    _ blockMode = iota
    repeat
    single
)

type data struct {
    curIndex int
    kIndex int
    vIndex int
    content string
    mode blockMode
}

type fieldInfo struct {
    value string
    isAddQuote bool
}

type CCondSqlSplice struct {
}

func (self *CCondSqlSplice) Serde(sql string, obj interface{}) (*string, error) {
    maps := map[int]*map[string]*[]*fieldInfo{}
    err := self.fields(obj, &maps)
    if err != nil {
        return nil, err
    }
    r := self.parse(sql, func(d *data) string {
        /*
        ** 比较 kIndex 与 vIndex 的大小
        ** 先替换较大的 (较小的索引就不需要变更)
        */
        switch d.mode {
        case repeat:
            var (
                v1, v2 int
                v1Len, v2Len int
                buf bytes.Buffer
            )
            if d.kIndex < d.vIndex {
                v1 = d.kIndex
                v2 = d.vIndex
                v1Len = keyword_key_len
                v2Len = keyword_value_len
            } else {
                v1 = d.vIndex
                v2 = d.kIndex
                v2Len = keyword_key_len
                v1Len = keyword_value_len
            }
            if v, ok := maps[d.curIndex]; ok {
                /*
                ** 遍历字段map
                */
                for key, va := range *v {
                    /*
                    ** 替换 v2
                    */
                    for _, value := range *va {
                        bufOnce := bytes.Buffer{}
                        bufOnce.WriteString(d.content[0:v2])
                        if value.isAddQuote {
                            bufOnce.WriteRune('\'')
                        }
                        bufOnce.WriteString(value.value)
                        if value.isAddQuote {
                            bufOnce.WriteRune('\'')
                        }
                        bufOnce.WriteString(d.content[v2+v2Len:])
                        /*
                        ** 替换 v1
                        */
                        t := bufOnce.String()
                        bufOnce.Reset()
                        bufOnce.WriteString(t[0:v1])
                        bufOnce.WriteString(key)
                        bufOnce.WriteString(t[v1+v1Len:])
                        buf.WriteString(bufOnce.String())
                    }
                }
            } else {
                return buf.String()
            }
            return buf.String()
        case single:
            var buf bytes.Buffer
            if v, ok := maps[d.curIndex]; ok {
                /*
                ** 遍历字段map
                */
                for _, va := range *v {
                    for _, value := range *va {
                        bufOnce := bytes.Buffer{}
                        bufOnce.WriteString(d.content[0:d.vIndex])
                        if value.isAddQuote {
                            bufOnce.WriteRune('\'')
                        }
                        bufOnce.WriteString(value.value)
                        if value.isAddQuote {
                            bufOnce.WriteRune('\'')
                        }
                        bufOnce.WriteString(d.content[d.vIndex+keyword_va_len:])
                        buf.WriteString(bufOnce.String())
                    }
                }
            }
            return buf.String()
        }
        return ""
    })
    return &r, nil
}

func (self *CCondSqlSplice) fields(obj interface{}, maps *map[int]*map[string]*[]*fieldInfo) error {
    var valueType reflect.Type
    value := reflect.ValueOf(obj)
    outType := reflect.TypeOf(obj)
    outputKind := outType.Kind()
    switch outputKind {
    case reflect.Ptr:
        valueType = outType.Elem()
        value = value.Elem()
    default:
        valueType = reflect.TypeOf(obj)
    }
    /*
    ** valueType: 指针指向的类型
    */
    switch valueType.Kind() {
    case reflect.Struct:
        self.obj2MapStrStrStructInner(value, valueType, maps)
    case reflect.Interface:
        self.obj2MapStrStrStructInner(value, valueType, maps)
    case reflect.Map:
        // obj2MapStrStrMapInner(value, valueType)
    default:
        return errors.New("is not parse")
    }
    return nil
}

func (self *CCondSqlSplice) posSplit(posStr string) []int {
    var (
        v int64
        err error
        poses []int
    )
    word := bytes.Buffer{}
    for _, c := range posStr {
        switch c {
        case ',':
            v, err = strconv.ParseInt(word.String(), 10, 32)
            if err == nil {
                poses = append(poses, int(v))
            }
            word.Reset()
        case ' ':
        default:
            word.WriteRune(c)
        }
    }
    v, err = strconv.ParseInt(word.String(), 10, 32)
    if err == nil {
        poses = append(poses, int(v))
    }
    return poses
}

func (self *CCondSqlSplice) getFieldValue(field reflect.Value, values *[]*fieldInfo) {
    fieldKind := field.Type().Kind()
    var fieldValue string
    var isAddQuote bool = false
    var selfIsAdd = true
    switch fieldKind {
    case reflect.Ptr:
        isNil := field.IsNil()
        if isNil {
            return
        }
        fieldPtrType := field.Type().Elem()
        fieldPtrKind := fieldPtrType.Kind()
        fieldElem := field.Elem()
        switch fieldPtrKind {
        case reflect.String:
            fieldValue = fieldElem.String()
            isAddQuote = true
        case reflect.Int:
            i := fieldElem.Int()
            fieldValue = strconv.FormatInt(i, 10)
        case reflect.Bool:
            b := fieldElem.Bool()
            fieldValue = strconv.FormatBool(b)
        case reflect.Float64:
            f := fieldElem.Float()
            fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
        case reflect.Slice:
            l := fieldElem.Len()
            for i := 0; i < l; i++ {
                idxValue := fieldElem.Index(i)
                self.getFieldValue(idxValue, values)
            }
            selfIsAdd = false
        }
    default:
        switch fieldKind {
        case reflect.String:
            fieldValue = field.String()
            isAddQuote = true
        case reflect.Int:
            i := field.Int()
            fieldValue = strconv.FormatInt(i, 10)
        case reflect.Bool:
            b := field.Bool()
            fieldValue = strconv.FormatBool(b)
        case reflect.Float64:
            f := field.Float()
            fieldValue = strconv.FormatFloat(f, 'f', 10, 32)
        case reflect.Slice:
            l := field.Len()
            for i := 0; i < l; i++ {
                idxValue := field.Index(i)
                self.getFieldValue(idxValue, values)
            }
            selfIsAdd = false
        case reflect.Interface:
            isNil := field.IsNil()
            if isNil {
                return
            }
            self.getFieldValue(field.Elem(), values)
            selfIsAdd = false
        }
    }
    if selfIsAdd {
        *values = append(*values, &fieldInfo{
            value: fieldValue,
            isAddQuote: isAddQuote,
        })
    }
}

func (self *CCondSqlSplice) obj2MapStrStrStructInner(value reflect.Value, valueType reflect.Type, maps *map[int]*map[string]*[]*fieldInfo) {
    fieldNum := value.NumField()
    for i := 0; i < fieldNum; i++ {
        field := value.Field(i)
        fieldType := valueType.Field(i)
        var fieldName string
        fieldTag := fieldType.Tag
        condField := fieldTag.Get(cond_tag_cond_field)
        if condField != "" {
            fieldName = condField
        } else {
            fieldName = fieldTag.Get(cond_tag_field)
        }
        if fieldName == "" {
            fieldName = fieldType.Name
        }
        posStr := fieldTag.Get(cond_tag_pos)
        if posStr == "" {
            continue
        }
        var isQuota bool = true
        quota := fieldTag.Get(cond_tag_quota)
        if quota != "" {
            switch quota {
            case keyword_quota_false:
                isQuota = false
            case keyword_quota_true:
                isQuota = true
            }
        }
        /*
        ** 如果 pos 解析为 int 失败, 则默认为0
        */
        poses := self.posSplit(posStr)
        if len(poses) == 0 {
            continue
        }
        for _, pos := range poses {
            values := []*fieldInfo{}
            self.getFieldValue(field, &values)
            /*
            ** 如果pos 已经存在, 则追加, 否则新建
            */
            for _, v := range values {
                fieldValue := v.value
                isAddQuote := v.isAddQuote
                if !isQuota {
                    isAddQuote = false
                }
                // fmt.Println(pos, fieldValue, isAddQuote, fieldName)
                if v, ok := (*maps)[int(pos)]; ok {
                    info := &fieldInfo{
                        value: fieldValue,
                        isAddQuote: isAddQuote,
                    }
                    if va, o := (*v)[fieldName]; o {
                        *va = append(*va, info)
                    } else {
                        a := []*fieldInfo{}
                        a = append(a, info)
                        (*v)[fieldName] = &a
                    }
                } else {
                    m := map[string]*[]*fieldInfo{}
                    a := []*fieldInfo{}
                    a = append(a, &fieldInfo{
                        value: fieldValue,
                        isAddQuote: isAddQuote,
                    })
                    m[fieldName] = &a
                    (*maps)[int(pos)] = &m
                }
            }
        }
    }
}

func (*CCondSqlSplice) parse(sql string, fn func(*data) string) string {
    buf := bytes.Buffer{}
    word := bytes.Buffer{}
    content := bytes.Buffer{}
    startIndex := 0
    kIndexTmp := 0
    curIndex := -1
    var mode Mode = normal
    for i, c := range sql {
        switch mode {
        case normal:
            switch c {
            case '{':
                mode = bigBrackets
                startIndex = i
                curIndex += 1
            default:
                buf.WriteRune(c)
            }
        case bigBrackets:
            switch c {
            case '}':
                switch word.String() {
                case keyword_key:
                    kIndexTmp = i - 1 - startIndex - keyword_key_len
                case keyword_value:
                    vIndex := i - 1 - startIndex - keyword_value_len
                    buf.WriteString(fn(&data{
                        curIndex: curIndex,
                        kIndex: kIndexTmp,
                        vIndex: vIndex,
                        content: content.String(),
                        mode: repeat,
                    }))
                case keyword_va:
                    vIndex := i - 1 - startIndex - keyword_va_len
                    buf.WriteString(fn(&data{
                        curIndex: curIndex,
                        vIndex: vIndex,
                        content: content.String(),
                        mode: single,
                    }))
                default:
                    // word.Reset()
                    // continue
                }
                mode = normal
                word.Reset()
                content.Reset()
            default:
                content.WriteRune(c)
                if c == ',' || c == '=' || c == ' ' || c == '.' || c == '(' || c == ')' {
                    switch word.String() {
                    case keyword_key:
                        kIndexTmp = i - 1 - startIndex - keyword_key_len
                    case keyword_value:
                        vIndex := i - 1 - startIndex - keyword_value_len
                        buf.WriteString(fn(&data{
                            curIndex: curIndex,
                            kIndex: kIndexTmp,
                            vIndex: vIndex,
                            content: content.String(),
                            mode: repeat,
                        }))
                    case keyword_va:
                        vIndex := i - 1 - startIndex - keyword_va_len
                        buf.WriteString(fn(&data{
                            curIndex: curIndex,
                            vIndex: vIndex,
                            content: content.String(),
                            mode: single,
                        }))
                    default:
                        word.Reset()
                        continue
                    }
                    word.Reset()
                } else {
                    word.WriteRune(c)
                }
            }
        }
    }
    return buf.String()
}

func NewCondSqlSplice() *CCondSqlSplice {
    s := CCondSqlSplice{}
    return  &s
}

func CondSqlSplice(sql string, obj interface{}) (*string, error) {
    s := CCondSqlSplice{}
    return s.Serde(sql, obj)
}
