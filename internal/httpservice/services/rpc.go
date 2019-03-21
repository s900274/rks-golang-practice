package services

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/downstream"
    logger "github.com/shengkehua/xlog4go"
)

func GetConnStats() map[string] define.ConnStats{
    var apiConnStats = make(map[string] define.ConnStats)

    for apiNmae, connObj := range downstream.DownStreamMgr{
        cliConnCnt := make(map[string]define.PoolStatus)
        connStats := define.ConnStats{}
        connStats.Addrs = make(map[string] define.AddrStats)
        connStats.NoInAddrs = make(map[string] define.PoolStatus)

        logger.Info("> GetConnStats, %s Start", apiNmae)
        //取得 這個ApiName的所有Conn
        for cli, _ := range connObj.MapConn {
            cliAddr, inPool, _, _ := connObj.GetMapConn(cli)
            p :=0
            u :=0
            switch inPool {
            case true:
                p ++
            default:
                u ++
            }
            if poolStruct, ok := cliConnCnt[cliAddr]; ok {
                poolStruct.InPool += p
                poolStruct.OnUse += u
                cliConnCnt[cliAddr] = poolStruct
            }else{
                _ps := define.PoolStatus{
                    InPool: p,
                    OnUse: u,
                }
                cliConnCnt[cliAddr] = _ps
            }
        }

        //整理 所有的Addr資訊
        for addr, _ := range connObj.GetAllAddr(){
            health, errCnt, _, connCnt := connObj.GetAddrAttr(addr)
            addrObj := define.AddrStats{
                Health: health,
                ErrCount: errCnt,
                ConnCount: connCnt,
            }

            if val, ok := cliConnCnt[addr]; ok{
                addrObj.InPoolCount = val.InPool
                addrObj.OnUseCount = val.OnUse
            }else{
                addrObj.InPoolCount = -1
                addrObj.OnUseCount = -1
            }
            connStats.Addrs[addr] = addrObj
        }

        //有在pool中, 但沒有出現在Addr的
        for cliAddr, v := range cliConnCnt {
            //找不存在
            if _, ok := connStats.Addrs[cliAddr]; !ok {
                connStats.NoInAddrs[cliAddr] = v
            }
        }

        apiConnStats[apiNmae] = connStats
    }
    return apiConnStats
}

func Dashboard(c *gin.Context) {
    result := GetConnStats()
    t := c.Query("type")
    if t == "json"{
        c.JSON(http.StatusOK, result)
    }else{
        c.HTML(http.StatusOK, "dashboard.tmpl", result)
    }

    return
}
