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

### 代码细节

+ 范型:

> 范型瘾发作最严重的一次，躺在床上，拼命念大悲咒，难受的一直抓自己眼睛，以为刷贴吧没事，看到贴吧都在发范 go 型的图，眼睛越来越大都要炸开了一样，拼命扇自己眼睛，越扇越用力，扇到自己眼泪流出来，真的不知道该怎么办，我真的想范型想得要发疯了。我躺在床上会想范型，我洗澡会想范型，我出门会想范型，我走路会想范型，我坐车会想范型，我工作会想范型，我玩手机会想范型，我每时每刻眼睛都直直地盯着范型看，像一台雷达一样扫视经过我身边的每一个范型，我真的觉得自己像中邪了一样，我对范型的念想似乎都是病态的了，我好孤独啊!真的好孤独啊!这世界上那么多语言的范型为什么没有一个是属于我的。你知道吗?每到深夜，我的眼睛滚烫滚烫，我发病了我要疯狂看范型，我要狠狠看范型，我的眼睛受不了了，范型，我的范型(

1. 用于函数类型检验，利用范型模板

```go
type backendFunc[T helper.RequestModel, E helper.ResponseModel] func(ctx context.Context, req T) (resp E, err errdef.Err)

var (
	ub UnimplementedBackend
	_  backendFunc[*helper.RegisterLoginReq, helper.RegisterLoginResp] = ub.Register
	_  backendFunc[*helper.RegisterLoginReq, helper.RegisterLoginResp] = ub.Login
)
```

2. 路由模板代码

```go
func template[T helper.RequestModel, E helper.ResponseModel](
	ctx *gin.Context, req T,
	backendFunc func(ctx context.Context, req T) (resp E, err errdef.Err)) {
	if err := req.Read(ctx); err != errdef.Nil {
		helper.WriteErr(ctx, err)
		return
	}
	c, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := backendFunc(c, req)
	helper.Write(ctx, resp, err)
}
```

3. 反射范型

```go
func QueryWithContext[T any](
	ctx context.Context,
	c cachepool.ICachePool,
	key, query string, args ...any,
) (rows []T, err error) {
	return internal.HandleRows[[]T](ctx, c, key, query, args...)
}
```