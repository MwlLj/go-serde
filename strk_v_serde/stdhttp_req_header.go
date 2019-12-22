package kv_serde

import (
    "net/http"
)

type httpHeader struct {
    req *http.Request
}

func (self *httpHeader) init(req *http.Request) error {
    self.req = req
    return nil
}

func (self *httpHeader) get(key *string) string {
    return self.req.Header.Get(*key)
}

func HttpReqHeaderDeserde(req *http.Request, output interface{}) error {
    input := httpHeader{}
    return deserde(&input, req, output)
}
