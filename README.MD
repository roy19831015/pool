# Pool

"Pool" is a module used for pooling every resource interface in golang


  <h3 align="center">Pool</h3>
  <p align="center">
    "Pool" is a module used for pooling every resource interface in golang.
  </p>


 This very README.md is for gopher.



### Requirements

golang 1.18 ……Because the feature 'Generics' is supported by golang after the version 1.18

### **Getting Started**

try go get as below or **just use it in go.mod**

```sh
go get github.com/roy19831015/pool/pool
```

### Example
The <a href="https://github.com/roy19831015/pool/blob/master/cmd/websocket_connection_pool_test_main.go">cmd/websocket_connection_pool_test_main.go</a> provides an example that used for pooling websocket connections.

In this example, main function tries to make connections between client and server, client side has 100 goroutines which are limited competing to use at most 10 connections.

It works well.

### Notice

I recomment that the 'Generics' [T any] should be an interface type.

In other words, comparing to 'Conn', '*Conn' is recommended.

### Copyright

You can use all the codes in this repo at the moment when you see this.