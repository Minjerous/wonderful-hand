# game

> game 在线服务是整个象棋游戏的服务器
> 
> 网络层是基于 websocket 的自定的协议(具体见 network/protocol)
> 

##  有以下功能

> + 维持客户端与服务端的 ws(tcp) 连接
> + 握手验证客户端
> + 接受和处理客户端的数据包 (数据包见 network/protocol/packet)
> + 提供和执行 rpc 服务
>   + 与 room 和 user 服务交互
>   + 提供 rpc 让其他服务可以与客户端交互

## 特色

> + epoll 特色的服务器，可以承载更多连接 (gnet)
> + RPC ，微服务化
> + 通信用的是二进制的数据包，这样效率更高