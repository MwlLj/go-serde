package kv_serde

import (
    "net/http"
    "net/url"
)

type httpParam struct {
    values url.Values
}

func (self *httpParam) init(req *http.Request) error {
    self.values = req.URL.Query()
    return nil
}

func (self *httpParam) get(key *string) string {
    return self.values.Get(*key)
}

func HttpReqParamDeserde(req *http.Request, output interface{}) error {
    input := httpParam{}
    return deserde(&input, req, output)
}

func HttpReqParamValue(req *http.Request, key string) *string {
    if v := req.URL.Query().Get(key); v != "" {
        return &v
    } else {
        return nil
    }
}
