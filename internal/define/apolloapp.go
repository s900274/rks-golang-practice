package define


type ApolloConf struct {
    AppId string
    Cluster string
    NamespaceName string
    Ip string
}

type ServiceInfo struct {
    Servers []string
    Protocol string
}

type DefaultConfiguration struct {
    ApiNameList []string
    ApiInfo     map[string]ServiceInfo
}