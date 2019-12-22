package kv_serde

import (
    "testing"
    "net/http"
    "fmt"
    "time"
)

type httpHeaderDeserdeExtra struct {
    V1 *string `json:"v1"`
}

type httpHeaderDeserdeStruct struct {
    Name *string `field:"name"`
    Age int `field:"age"`
    Extra *httpHeaderDeserdeExtra `field:"extra"`
}

func TestHttpHeaderDeserde(t *testing.T) {
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            p := httpParamDeserdeStruct{}
            if err := HttpReqHeaderDeserde(r, &p); err != nil {
                fmt.Println(err)
            }
            fmt.Println(*p.Name, p.Age, *p.Extra.V1)
        })
        http.ListenAndServe(":50001", mux)
    }()
    <-time.After(1 * time.Second)
    fmt.Println("send request")
    url := `http://127.0.0.1:50001/`
    method := "GET"
    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    req.Header.Add("name", "Mike")
    req.Header.Add("age", "25")
    req.Header.Add("extra", `{"v1":"456"}`)
    res, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer res.Body.Close()
}
