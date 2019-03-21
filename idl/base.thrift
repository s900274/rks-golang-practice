namespace cpp bc.base
namespace go bc.base
namespace java com.magic.base
namespace php bc.base

struct Trace {
    ///< 调用方logid
    1: required string logId;
    ///< 调用方系统名
    2: required string caller;
}

// signature加密算法：
// MD5(keystr:MD5(version=$version:type=$type:cmd=$cmd:timestamp=$timestamp:body=$body))
// MD5加密后的字符串全部使用大写字母


struct MarsHeader {
    ///< 协议版本号
    1: required byte    version;
    ///< ms毫秒数
    2: required i64     timestamp;
    ///< 验证签名
    3: required string  signature;
    ///< 用户token
    4: optional string  token;
}

struct MarsRequest {
    ///< 服務名稱
    1: required string  apiname;
    ///< 命令識別
    2: required i64 action;
    ///< 请求协议头
    3: required MarsHeader header;
    ///< 请求协议体：JSON格式
    4: required string   body;
}

struct MarsResponse {
    ///< 狀態碼
    1: required i32     code;
    ///< Msg
    2: required string  msg;
    ///< 返回
    3: optional string  data;
}