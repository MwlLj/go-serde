package kv_serde

import (
    "testing"
    "net/http"
    "fmt"
    "time"
)

type httpParamDeserdeExtra struct {
    V1 *string `json:"v1"`
}

type httpParamDeserdeStruct struct {
    Name *string `field:"name"`
    Age int `field:"age"`
    Extra *httpParamDeserdeExtra `field:"extra"`
}

func TestHttpParamDeserde(t *testing.T) {
    go func() {
        mux := http.NewServeMux()
        mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            p := httpParamDeserdeStruct{}
            if err := HttpReqParamDeserde(r, &p); err != nil {
                fmt.Println(err)
            }
            fmt.Println(*p.Name, p.Age, *p.Extra.V1)
        })
        http.ListenAndServe(":50000", mux)
    }()
    <-time.After(1 * time.Second)
    fmt.Println("send request")
    http.Get(`http://127.0.0.1:50000/?name=Jake&age=20&extra={"v1":"123"}`)
}
