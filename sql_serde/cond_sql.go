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
    cond_tag_pos string = "pos"
    keyword_key string = "k"
    keyword_value string = "v"
)

var (
    keyword_key_len int = len(keyword_key)
    keyword_value_len int = len(keyword_value)
)

type Mode int8
const (
    _ Mode = iota
    normal
    bigBrackets
)

type data struct {
    curIndex int
    kIndex int
    vIndex int
    content string
}

type fieldInfo struct {
    value string
    isAddQuote bool
}

type CCondSqlSplice struct {
}

func (self *CCondSqlSplice) Serde(sql string, obj interface{}) (*string, error) {
    maps := map[int]*map[string]*fieldInfo{}
    err := self.fields(obj, &maps)
    if err != nil {
        return nil, err
    }
    r := self.parse(sql, func(d *data) string {
        /*
        ** 比较 kIndex 与 vIndex 的大小
        ** 先替换较大的 (较小的索引就不需要变更)
        */
        var (
            v1, v2 int
            buf bytes.Buffer
        )
        if d.kIndex < d.vIndex {
            v1 = d.kIndex
            v2 = d.vIndex
        } else {
            v1 = d.vIndex
            v2 = d.kIndex
        }
        if v, ok := maps[d.curIndex]; ok {
            /*
            ** 遍历字段map
            */
            for key, value := range *v {
                /*
                ** 替换 v2
                */
                bufOnce := bytes.Buffer{}
                bufOnce.WriteString(d.content[0:v2])
                if value.isAddQuote {
                    bufOnce.WriteRune('"')
                }
                bufOnce.WriteString(value.value)
                if value.isAddQuote {
                    bufOnce.WriteRune('"')
                }
                bufOnce.WriteString(d.content[v2+1:])
                /*
                ** 替换 v1
                */
                t := bufOnce.String()
                bufOnce.Reset()
                bufOnce.WriteString(t[0:v1])
                bufOnce.WriteString(key)
                bufOnce.WriteString(t[v1+1:])
                buf.WriteString(bufOnce.String())
            }
        } else {
            return buf.String()
        }
        return buf.String()
    })
    return &r, nil
}

func (self *CCondSqlSplice) fields(obj interface{}, maps *map[int]*map[string]*fieldInfo) error {
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

func (self *CCondSqlSplice) obj2MapStrStrStructInner(value reflect.Value, valueType reflect.Type, maps *map[int]*map[string]*fieldInfo) {
    fieldNum := value.NumField()
    for i := 0; i < fieldNum; i++ {
        field := value.Field(i)
        fieldType := valueType.Field(i)
        var fieldName string
        fieldTag := fieldType.Tag
        fieldName = fieldTag.Get(cond_tag_field)
        if fieldName == "" {
            fieldName = fieldType.Name
        }
        posStr := fieldTag.Get(cond_tag_pos)
        if posStr == "" {
            continue
        }
        /*
        ** 如果 pos 解析为 int 失败, 则默认为0
        */
        pos, err := strconv.ParseInt(posStr, 10, 32)
        if err != nil {
            continue
        }
        fieldKind := field.Type().Kind()
        var fieldValue string
        var isAddQuote bool = false
        switch fieldKind {
        case reflect.Ptr:
            isNil := field.IsNil()
            if isNil {
                continue
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
            }
        }
        /*
        ** 如果pos 已经存在, 则追加, 否则新建
        */
        if v, ok := (*maps)[int(pos)]; ok {
            (*v)[fieldName] = &fieldInfo{
                value: fieldValue,
                isAddQuote: isAddQuote,
            }
        } else {
            m := map[string]*fieldInfo{}
            m[fieldName] = &fieldInfo{
                value: fieldValue,
                isAddQuote: isAddQuote,
            }
            (*maps)[int(pos)] = &m
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
                    kIndexTmp = i - 1 - startIndex - 1
                case keyword_value:
                    vIndex := i - 1 - startIndex - 1
                    buf.WriteString(fn(&data{
                        curIndex: curIndex,
                        kIndex: kIndexTmp,
                        vIndex: vIndex,
                        content: content.String(),
                    }))
                default:
                    word.Reset()
                    continue
                }
                mode = normal
                word.Reset()
                content.Reset()
            default:
                content.WriteRune(c)
                if c == ',' || c == '=' || c == ' ' {
                    switch word.String() {
                    case keyword_key:
                        kIndexTmp = i - 1 - startIndex - 1
                    case keyword_value:
                        vIndex := i - 1 - startIndex - 1
                        buf.WriteString(fn(&data{
                            curIndex: curIndex,
                            kIndex: kIndexTmp,
                            vIndex: vIndex,
                            content: content.String(),
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
