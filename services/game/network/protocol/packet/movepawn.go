package packet

import (
	"wonderful-hand-game/network/protocol"
)

type PawnID = uint8

const (
	IDPawnRedBoss        PawnID = iota // 红帅
	IDPawnRedCarL                      // 红车左
	IDPawnRedCarR                      // 红车右
	IDPawnRedHorseL                    // 红马左
	IDPawnRedHorseR                    // 红马右
	IDPawnRedElephantL                 // 红相左
	IDPawnRedElephantR                 // 红相右
	IDPawnRedKnightL                   // 红仕左
	IDPawnRedKnightR                   // 红仕右
	IDPawnRedGunL                      // 红炮左
	IDPawnRedGunR                      // 红炮右
	IDPawnRedSoldier1                  // 红兵1 (从左向右)
	IDPawnRedSoldier2                  // 红兵2 (从左向右)
	IDPawnRedSoldier3                  // 红兵3 (从左向右)
	IDPawnRedSoldier4                  // 红兵4 (从左向右)
	IDPawnRedSoldier5                  // 红兵5 (从左向右)
	_                                  // 16
	IDPawnBlackBoss                    // 黑帅
	IDPawnBlackCarL                    // 黑车左
	IDPawnBlackCarR                    // 黑车右
	IDPawnBlackHorseL                  // 黑马左
	IDPawnBlackHorseR                  // 黑马右
	IDPawnBlackElephantL               // 黑象左
	IDPawnBlackElephantR               // 黑象右
	IDPawnBlackKnightL                 // 黑士左
	IDPawnBlackKnightR                 // 黑士右
	IDPawnBlackGunL                    // 黑炮左
	IDPawnBlackGunR                    // 黑炮右
	IDPawnBlackSoldier1                // 黑卒1 (从左向右)
	IDPawnBlackSoldier2                // 黑卒2 (从左向右)
	IDPawnBlackSoldier3                // 黑卒3 (从左向右)
	IDPawnBlackSoldier4                // 黑卒4 (从左向右)
	IDPawnBlackSoldier5                // 黑卒5 (从左向右)
)

// MovePawn 由客户端发出告诉服务器自己移动的棋子
// 由服务端发出同步给对方棋子的移动
type MovePawn struct {
	PawnID         PawnID
	UID            uint64
	Token          string
	RoomID         string
	DeltaX, DeltaY uint8
}

func (m *MovePawn) ID() uint8 {
	return IDMovePawn
}

func (m *MovePawn) Marshal(w *protocol.Writer) {
	w.Uint8(&m.PawnID)
	w.VarUint64(&m.UID)
	w.String(&m.Token)
	w.String(&m.RoomID)
	w.Uint8(&m.DeltaX)
	w.Uint8(&m.DeltaY)
}

func (m *MovePawn) Unmarshal(r *protocol.Reader) {
	r.Uint8(&m.PawnID)
	r.VarUint64(&m.UID)
	r.String(&m.Token)
	r.String(&m.RoomID)
	r.Uint8(&m.DeltaX)
	r.Uint8(&m.DeltaY)
}
