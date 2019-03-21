package define


type ConnStats struct {
    Addrs map[string] AddrStats
    NoInAddrs map[string] PoolStatus
}

type AddrStats struct {
    Health bool
    ErrCount int64
    ConnCount int
    InPoolCount int
    OnUseCount int
}

type PoolStatus struct {
    InPool int
    OnUse int
}