package gorpc

import (
    "time"
)

type Server struct {
    opts *ServerOptions
    services map[string]Service
}

type ServerOptions struct {
    address string              // listening address, e.g. :( ip://127.0.0.1:8080、 dns://www.google.com)
    network string              // network type, e.g. : tcp、udp
    protocol string             // protocol typpe, e.g. : proto、json
    timeout time.Duration       // timeout
    serializationType string    // serialization type, default: proto
 
    selectorSvrAddr string       // service discovery server address, required when using the third-party service discovery plugin
    tracingSvrAddr  string         // tracing plugin server address, required when using the third-party tracing plugin
    tracingSpanName string       // tracing span name, required when using the third-party tracing plugin
    pluginNames []string         // plugin name
    interceptors []interceptor.ServerInterceptor
}

// Service 定义了某个具体服务的通用实现接口
type Service interface {
    Register(string, Handler)
    Serve(*ServerOptions)
    Close()
}

type service struct{
    svr interface{}          // server
    ctx context.Context       // 每一个 service 一个上下文进行管理
    cancel context.CancelFunc   // context 的控制器
    serviceName string        // 服务名
    handlers map[string]Handler
    opts *ServerOptions       // 参数选项
 }
