# grpc-test
grpc-test 练习

在proto目录执行: 

protoc --go_out=plugins=grpc:. 1.proto 

生成pb.go文件

分别有四种调用方式

1.rpc请求 请求的函数 (发送请求参数) returns (返回响应的参数)

2.rpc请求 客户端流式(请求参数+stream标识)客户端传入多个请求对象，服务端返回一个响应结果。

3.rpc请求 服务端流式(返回参数+stream标识)一个请求对象,服务端返回多个结果对象。

4.rpc请求 服务端客户端双流(请求+返回参数都需要 +stream标识)传入多个对象可以返回多个响应对象
