package downstream

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/helpers/ucm"
)

var DownStreamMgr = make(map[string]*ucm.ChanConnPool)


func Register(name string, inst *ucm.ChanConnPool) {
    if inst == nil {
        panic("downstream: Register adapter is nil")
    }
    if _, ok := DownStreamMgr[name]; ok {
        panic("downstream: Register called twice for adapter " + name)
    }
    DownStreamMgr[name] = inst
}

