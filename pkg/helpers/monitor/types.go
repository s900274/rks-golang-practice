package monitor
//
//type Count struct {
//
//	// gomonitor chan的数据量
//	MonitorChanSize int64 `json:"MonitorChanSize"`
//
//    Request_Config_12345 int64 `json:"Request_Config_12345"`
//    //
//	////  stg key 写codis次数
//	//STG_DRIVER_EMPTY_CNT                         int64 `json:"STG_DRIVER_EMPTY_CNT" type:"qps"`
//	//STG_DRIVER_TOTAL_CNT                         int64 `json:"STG_DRIVER_TOTAL_CNT" type:"qps"`
//	//STG_ORDER_CREATE_CNT                         int64 `json:"STG_ORDER_CREATE_CNT" type:"qps"`
//	//STG_ORDER_NOT_BROADCAST_CNT                  int64 `json:"STG_ORDER_NOT_BROADCAST_CNT" type:"qps"`
//	//STG_DRIVER_UNASSIGNED_TOTAL_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_TOTAL_CNT" type:"qps"`
//	//STG_DRIVER_UNASSIGNED_EMPTY_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_EMPTY_CNT" type:"qps"`
//	//STG_REALTIME_ORDER_CREATE_CNT                int64 `json:"STG_REALTIME_ORDER_CREATE_CNT" type:"qps"`
//	//STG_REALTIME_ORDER_NOT_BROADCAST_CNT         int64 `json:"STG_REALTIME_ORDER_NOT_BROADCAST_CNT" type:"qps"`
//	//STG_REALTIME_CARPOOL_ORDER_CREATE_CNT        int64 `json:"STG_REALTIME_CARPOOL_ORDER_CREATE_CNT" type:"qps"`
//	//STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_CNT int64 `json:"STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_CNT" type:"qps"`
//}
//
//type TimeUsed struct {
//	////  stg key 写codis平均时间
//	//STG_DRIVER_EMPTY_CNT                         int64 `json:"STG_DRIVER_EMPTY_AVGTime" type:"avg"`
//	//STG_DRIVER_TOTAL_CNT                         int64 `json:"STG_DRIVER_TOTAL_AVGTime" type:"avg"`
//	//STG_ORDER_CREATE_CNT                         int64 `json:"STG_ORDER_CREATE_AVGTime" type:"avg"`
//	//STG_ORDER_NOT_BROADCAST_CNT                  int64 `json:"STG_ORDER_NOT_BROADCAST_AVGTime" type:"avg"`
//	//STG_DRIVER_UNASSIGNED_TOTAL_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_TOTAL_AVGTime" type:"avg"`
//	//STG_DRIVER_UNASSIGNED_EMPTY_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_EMPTY_AVGTime" type:"avg"`
//	//STG_REALTIME_ORDER_CREATE_CNT                int64 `json:"STG_REALTIME_ORDER_CREATE_AVGTime" type:"avg"`
//	//STG_REALTIME_ORDER_NOT_BROADCAST_CNT         int64 `json:"STG_REALTIME_ORDER_NOT_BROADCAST_AVGTime" type:"avg"`
//	//STG_REALTIME_CARPOOL_ORDER_CREATE_CNT        int64 `json:"STG_REALTIME_CARPOOL_ORDER_CREATE_AVGTime" type:"avg"`
//	//STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_CNT int64 `json:"STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_AVGTime" type:"avg"`
//    //
//	////  stg key 写codis最大耗时
//	//MAX_STG_DRIVER_EMPTY_CNT                         int64 `json:"STG_DRIVER_EMPTY_MAXTime"`
//	//MAX_STG_DRIVER_TOTAL_CNT                         int64 `json:"STG_DRIVER_TOTAL_MAXTime"`
//	//MAX_STG_ORDER_CREATE_CNT                         int64 `json:"STG_ORDER_CREATE_MAXTime"`
//	//MAX_STG_ORDER_NOT_BROADCAST_CNT                  int64 `json:"STG_ORDER_NOT_BROADCAST_MAXTime"`
//	//MAX_STG_DRIVER_UNASSIGNED_TOTAL_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_TOTAL_MAXTime"`
//	//MAX_STG_DRIVER_UNASSIGNED_EMPTY_CNT              int64 `json:"STG_DRIVER_UNASSIGNED_EMPTY_MAXTime"`
//	//MAX_STG_REALTIME_ORDER_CREATE_CNT                int64 `json:"STG_REALTIME_ORDER_CREATE_MAXTime"`
//	//MAX_STG_REALTIME_ORDER_NOT_BROADCAST_CNT         int64 `json:"STG_REALTIME_ORDER_NOT_BROADCAST_MAXTime"`
//	//MAX_STG_REALTIME_CARPOOL_ORDER_CREATE_CNT        int64 `json:"STG_REALTIME_CARPOOL_ORDER_CREATE_MAXTime"`
//	//MAX_STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_CNT int64 `json:"STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_MAXTime"`
//}
//
//type Resource struct {
//	Mem int64   `json:"mem"`
//	Cpu float64 `json:"cpu"`
//
//	Goroutines int64 `json:"Goroutines"`
//	Fds        int64 `json:"Fds"`
//
//	Mem_Allocated int64 `json:"Mem_Allocated"`
//	Mem_Mallocs   int64 `json:"Mem_Mallocs"`
//	Mem_Heap      int64 `json:"Mem_Heap"`
//	Mem_Stack     int64 `json:"Mem_Stack"`
//	Mem_Objects   int64 `json:"Mem_Objects"`
//
//	Gc_Num   int64 `json:"Gc_Num"`
//	Gc_Pause int64 `json:"Gc_Pause"`
//	Gc_Next  int64 `json:"Gc_Next"`
//	Cgo      int64 `json:"Cgo"`
//}
//
//type Error struct {
//    //
//	////  stg key 写codis次数
//	//STG_DRIVER_EMPTY_CNT_ERRCNT                         int64 `json:"STG_DRIVER_EMPTY_ERRCNT"`
//	//STG_DRIVER_TOTAL_CNT_ERRCNT                         int64 `json:"STG_DRIVER_TOTAL_ERRCNT"`
//	//STG_ORDER_CREATE_CNT_ERRCNT                         int64 `json:"STG_ORDER_CREATE_ERRCNT"`
//	//STG_ORDER_NOT_BROADCAST_CNT_ERRCNT                  int64 `json:"STG_ORDER_NOT_BROADCAST_ERRCNT"`
//	//STG_DRIVER_UNASSIGNED_TOTAL_CNT_ERRCNT              int64 `json:"STG_DRIVER_UNASSIGNED_TOTAL_ERRCNT"`
//	//STG_DRIVER_UNASSIGNED_EMPTY_CNT_ERRCNT              int64 `json:"STG_DRIVER_UNASSIGNED_EMPTY_ERRCNT"`
//	//STG_REALTIME_ORDER_CREATE_CNT_ERRCNT                int64 `json:"STG_REALTIME_ORDER_CREATE_ERRCNT"`
//	//STG_REALTIME_ORDER_NOT_BROADCAST_CNT_ERRCNT         int64 `json:"STG_REALTIME_ORDER_NOT_BROADCAST_ERRCNT"`
//	//STG_DRIVER_MINUTE_CNT_ERRCNT                        int64 `json:"STG_DRIVER_MINUTE_ERRCNT"`
//	//STG_REALTIME_CARPOOL_ORDER_CREATE_CNT_ERRCNT        int64 `json:"STG_REALTIME_CARPOOL_ORDER_CREATE_ERRCNT"`
//	//STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_CNT_ERRCNT int64 `json:"STG_REALTIME_CARPOOL_ORDER_NOT_BROADCAST_ERRCNT"`
//}


type Data struct {
	Count    map[string]int64 `json:"count"`
	Timeused map[string]int64 `json:"timeused"`
    TimeMax  map[string]int64 `json:"timemax"`
	Resource map[string]int64 `json:"resource"`
	Err      map[string]int64 `json:"error"`
}

type StatData struct {
    Type    string
    Cmd     string
    Cnt     int64
    Avg     int64
    Max     int64
    Err     int64
}
