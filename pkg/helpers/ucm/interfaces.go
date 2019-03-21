/*
通用连接管理器接口
20160417
@aosen
*/
package ucm

type Ucmer interface {
	//获取连接
	Get() (interface{}, error)
	//将连接放回池中, 如果是短连接就是关闭
	Put(client interface{}, safe bool)
	//获取当前活跃连接数
	ActiveCount() int
	//注：client不能为nil，必须是由Get生成的合法连接
	//其实这个方法应该只用在连接池的情况下，短连接如果
	//生成连接失败。client=nil
	//应用场景：当Get到的连接被调用方close后，需要重连的情况下
	ReConn(client interface{}) (interface{}, error)
	//获取地址的健康状态
	GetHealthy() map[string]bool
	//获取地址对应的总连接数字典
	GetConnCount() map[string]int
}
