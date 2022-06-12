# game-protocol

> 这是整个游戏的数据包协议

## 数据结构

| Name      | Size                  | Type   |
| --------- | --------------------- | ------ |
| Magic     | 8 btyes               | []byte |
| VarUint32 | 4 bytes               | uint32 |
| ByteSlice | 4 + len(slice) bytes  | []byte |
| String    | 4 + len(string) bytes | string |
| Bool      | 1 byte                | bool   |
| Uint8     | 1 byte                | uint8  |
| Int8      | 1 byte                | int8   |
| VarUint64 | 8 btyes               | uint64 |

## 数据包帧结构

### TextPacket

> 聊天文本数据包

| PacketID | Bound To        |
| -------- | --------------- |
| 0x01     | Server & Client |

| FieldName  | FieldType | Notes                                                |
| ---------- | --------- | ---------------------------------------------------- |
| TextType   | Uint8     | 文本类型                                             |
| SourceName | String    | 源名称                                               |
| DestRoomID | String    | 目标房间ID (TextType 为 Global 无需定义)             |
| Content    | String    | 内容                                                 |
| UID        | VarUint64 | 当 TextType 为 1(TextTypeChat)时定义，源用户的 UID   |
| Token      | String    | 当 TextType 为 1(TextTypeChat)时定义，源用户的 Token |

### UnConnectRoomPingPacket

> 未连接的 Room Ping packet

| PacketID | Bound To |
| -------- | -------- |
| 0x02     | Client   |

| FieldName  | FieldType | Notes      |
| ---------- | --------- | ---------- |
| Magic      | Magic     |            |
| DestRoomID | String    | 目标房间ID |

### UnConnectRoomPongPacket

> UnConnectRoomPingPacket 回应

| PacketID | Bound To |
| -------- | -------- |
| 0x03     | Server   |

| FieldName    | FieldType | Notes          |
| ------------ | --------- | -------------- |
| Magic        | Magic     |                |
| RoomName     | String    | 目标房间名称   |
| RoomSubtitle | String    | 目标房间小标题 |
| RoomStatus   | Uint8     | 房间状态       |

### OpenConnect1Packet

> 客户端发起连接请求

| PacketID | Bound To |
| -------- | -------- |
| 0x04     | Client   |

| FieldName | FieldType | Notes      |
| --------- | --------- | ---------- |
| Magic     | Magic     |            |
| UID       | VarUint64 | 用户 UID   |
| Name      | String    | 用户昵称   |
| Token     | String    | 用户 Token |

### OpenConnect2Packet

> 服务端响应连接请求

| PacketID | Bound To |
| -------- | -------- |
| 0x05     | Server   |

| FieldName  | FieldType | Notes  |
| ---------- | --------- | ------ |
| Magic      | Magic     |        |
| StatusCode | Uint8     | 状态码 |

### MovePawnPacket

> 移动棋子的数据包
>
> 该数据包由服务端发送会省略带有 * 的字段

| PacketID | Bound To        |
| -------- | --------------- |
| 0x06     | Server & Client |

| FieldName | FieldType | Notes         |
| --------- | --------- | ------------- |
| PawnID    | Uint8     | 棋子 ID       |
| UID*      | VarUint64 | 用户 UID      |
| Token*    | String    | 用户 Token    |
| RoomID*   | String    | 目标 Room ID  |
| DeltaX    | Uint8     | 棋子 x 偏移量 |
| DeltaY    | Uint8     | 棋子 y 偏移量 |
| Ack       | Uint8     | 见 Ack        |

> Ack: 由服务器发送 Ack 置 1 表示确认移动有效并广播给棋局玩家，Ack 置 2 表示无效返回给发送方；由客户端发送 Ack 置 0

### AddRoomPacket

> 添加房间的 Packet
>
> 服务器返回会省略带 * 字段

| PacketID | Bound To      |
| -------- | ------------- |
| 0x07     | Client&Server |

| FieldName | FieldType | Notes                |
| --------- | --------- | -------------------- |
| UID*      | VarUint64 | 用户 UID             |
| Token*    | String    | 用户 Token           |
| Password  | String    | 房间密码             |
| RoomID    | String    | 由服务器返回 房间 ID |

### JoinRoomPacket

> 加入房间的 Packet

| PacketID | Bound To      |
| -------- | ------------- |
| 0x08     | Client&Server |

| FieldName | FieldType | Notes      |
| --------- | --------- | ---------- |
| UID       | VarUint64 | 用户 UID   |
| Token     | String    | 用户 Token |
| RoomID    | String    | 房间 ID    |
| Password  | String    | 房间密码   |

### GameClaimPakcet

> 游戏确认 Packet
>
> 由客户端发送确认游戏开始，只有所有棋局玩家确认了才会开始游戏

| PacketID | Bound To |
| -------- | -------- |
| 0x09     | Client   |

| FieldName | FieldType | Notes      |
| --------- | --------- | ---------- |
| UID       | VarUint64 | 用户 UID   |
| Token     | String    | 用户 Token |
| RoomID    | String    | 房间 ID    |

### GameStartPacket

> 游戏开始 Packet

| PacketID | Bound To |
| -------- | -------- |
| 0x0A     | Server   |

| FieldName | FieldType | Notes                |
| --------- | --------- | -------------------- |
| Map       | ByteSlice | 棋盘初始布局         |
| Label     | Uint8     | 客户端在游戏中的一方 |

### GameEndPacket

> 棋局结束 Packet
>
> 由服务端发送提示客户端游戏已经结束

| PacketID | Bound To |
| -------- | -------- |
| 0x0B     | Server   |

| FieldName | FieldType | Notes    |
| --------- | --------- | -------- |
| RoomID    | String    | 房间 ID  |
| EndStatus | Uint8     | 结束状态 |