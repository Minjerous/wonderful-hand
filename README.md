# wonderful-hand —— 妙手

> 来自一个俗手写的下棋服务
> 
> 属实是无从下手
> 

> 很遗憾，没有写完 (目前只有 user 服务能够部署)
> 
> 目前整体是微服务架构，分为以下服务
> 
> + chess
> + game
> + room
> + user
> 

> 完成内容情况
> 
> + user 服务
> + game 网络层
> + 服务发现
> + 负载均衡
> 

> 预计完成内容
>
> + chess 服务模拟棋盘逻辑
> + game 网络层完善
> + room 房间服务
> + 实现游戏协议的客户端

## 详细内容

> 各个服务的 README
> 
> + [chess](./services/chess/README.md)
> + [game](./services/game/README.md)
> + [room](./services/room/README.md)
> + [user](./services/user/README.md)

### 技术细节

+ 微服务架构，在线游戏服耦合度低，RPC
+ epoll, eventloops 风格的在线游戏服，性能更高，承载连接数多
+ 自定游戏数据包协议，实现快速序列化和反序列化
+ Websocket 长连接
+ etcd 服务发现
+ 负载均衡
+ 二级缓存(本地和全局redis)维护 room，chess session