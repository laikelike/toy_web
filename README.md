# toy_web
简易版Go web框架
v1: 创建并启动了一个简单的web服务
v2: 对注册路由和服务启动进行抽象封装
v3: 对每一次http请求封装，使用context代表一次http请求；参考restful使用map的key将http method和path绑定，value是context，实现handler。使得http method决定了操作，http path决定了操作对象。
v4: 利用组合把http.handler接口和我自己实现的HandlerBasedOnMap接口组合在一起，使得sdkHttpServer结构体可以使用他们所有的方法，达到解耦的目的，方便此结构体增减方法，进一步抽象注册路由函数(Route),
v6: 引入AOP的概念，使用责任链（filter）实现AOP
v7:自己生成一个简单的路由树
v8:优化路由树：实现通配符匹配。
