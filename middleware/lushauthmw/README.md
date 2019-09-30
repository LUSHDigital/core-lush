# Auth Middleware
The package `core-lush/middleware/lushauthmw` is used to attach authentication information to requests and responses for REST and gRPC. To learn more about how to use auth inside of your application you should read the [documentation for the **core-lush/lushauth** package](https://github.com/LUSHDigital/core-lush/tree/master/lushauth#auth).

## Examples

### Attach gRPC auth middlewares to server

```go
server := grpc.NewServer(
    middleware.WithUnaryServerChain(
        lushauthmw.NewUnaryServerInterceptor(broker),
    ),
)
```
