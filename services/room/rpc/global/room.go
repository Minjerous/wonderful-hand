package global

type Room struct {
	RoomID     string
	OwnerUID   uint64
	Password   string // 空说明没有密码
	RoomStatus RoomStatus
}

const (
	StatusRunning uint8 = iota
	StatusENDING
	StatusWAITING
)

type RoomStatus struct {
	StatusCode uint8 // RUNNING(游戏中) / ENDING(结算中) / WAITING(等待中)
}
