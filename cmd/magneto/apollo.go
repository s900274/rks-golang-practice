package main

import (
    "gitlab.kingbay-tech.com/engine-lottery/magneto/pkg/downstream"
    "gitlab.kingbay-tech.com/engine-lottery/magneto/internal/define"
    logger "github.com/shengkehua/xlog4go"
)


func inSliceIface(v string, sl []string) bool {
    for _, vv := range sl {
        if vv == v {
            return true
        }
    }
    return false
}

func sliceIntersect(slice1, slice2 []string) (diffslice []string) {
    for _, v := range slice1 {
        if inSliceIface(v, slice2) {
            diffslice = append(diffslice, v)
        }
    }
    return
}


func findChangeType(chgServers []string, poolServers []string) ([]string, []string){
    addS := []string{}
    delS := []string{}

    innerAddr := sliceIntersect(chgServers, poolServers)

    //在更換清單, 且不在pool, 則為新增
    for _, v := range chgServers{
        if !(inSliceIface(v, innerAddr)){
            addS = append(addS, v)
        }
    }

    //在pool清單, 且不在更換清單, 則為刪除
    for _, v := range poolServers{
        if !(inSliceIface(v, innerAddr)){
            delS = append(delS, v)
        }
    }
    return addS, delS
}

func ChangeApiServiceFromApollo(apiName string, chgInfo *define.ServiceInfo) {
    logger.Info("change apollo, apiName: %s start", apiName)
    stream := define.Cfg.Mars
    if ds, ok :=downstream.DownStreamMgr[apiName]; ok {
        connAddrs := ds.Getal()
        logger.Info("> change apollo, apiName: %s, origAddrs: %v", apiName, connAddrs)

        addServers, delServers := findChangeType(chgInfo.Servers, connAddrs)
        logger.Info("> change apollo, apiName: %s, addServers: %v", apiName, addServers)
        logger.Info("> change apollo, apiName: %s, delServers: %v", apiName, delServers)

        //刪除不要的addr連線
        for _, v := range delServers{
            ds.DelAddr(v)
        }

        //新增addr連線
        for _, v := range addServers{
            ds.CreateServerIntoDownStream(v, stream.Connsize)
        }

        //把連線移除pool
        if len(delServers) >0 {
            ds.DropOutAddrFromIdle(delServers)
        }
        logger.Info("> change apollo, apiName: %s, newAddrs: %v", apiName, ds.Getal())
    }

    return
}
