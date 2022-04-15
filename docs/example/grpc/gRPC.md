## Install
* Protocol buffer 编译器 protoc
* Go plugins 插件 protoc-gen-go 和 protoc-gen-go-grpc

```shell
# Move helloworld.proto to ROOT/grpc
# cd ROOT then execute in root path
protoc --go_out=./ --go-grpc_out=./ ./grpc/helloworld.proto
```


## Docs 
https://grpc.io/docs/languages/go/quickstart  
https://github.com/grpc/grpc-go/tree/master/examples  
