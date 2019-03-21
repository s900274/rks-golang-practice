package server

import (
    "testing"
    "golang.org/x/text/encoding"
    "fmt"
)


//
//func TestSignature(t *testing.T) {
//
//
//    s := &UranusServer{}
//
//    req := base.NewEGReq()
//    req.Header = base.NewEGHeader()
//    req.Header.Version = 1
//    req.Header.Type = base.CmdType_PASSPORT
//    req.Header.Cmd = 1048579
//    req.Header.Timestamp = 1495011190190
//    req.Header.Signature = "B43A79E52A01F23491AE870E97147E13"
//    req.Body = "{\"userId\":109467,\"username\":\"jackson\"}"
//    s.CheckSignature(req)
//    return
//}

func TestDecode(t *testing.T) {
    dec := encoding.Decoder{}
    s := "哈哈"
    s1, err := dec.String(s)
    if err != nil {
        fmt.Printf(err.Error())
    }
    fmt.Printf("%x--%x", s, s1)
}