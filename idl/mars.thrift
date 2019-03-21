namespace cpp bc.engine.mars
namespace go bc.engine.mars
namespace java com.magic.mars
namespace php bc.engine.mars

include "base.thrift"


service MarsService {
    ///< 调用引擎网关接口
    base.MarsResponse PublicInterface(
    1: base.MarsRequest req,
    2: base.Trace trace,
    );
}

