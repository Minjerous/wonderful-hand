package protocol

import (
	"io"
	"unsafe"
)

// Writer 将内建的数据结构写入 buffer 内
type Writer struct {
	w interface {
		io.Writer
		io.ByteWriter
	}
}

func NewWriter(w interface {
	io.Writer
	io.ByteWriter
}) *Writer {
	return &Writer{w: w}
}

func (w *Writer) Magic() {
	_, _ = w.w.Write(Magic)
}

func (w *Writer) VarUint32(x *uint32) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}

func (w *Writer) ByteSlice(x *[]byte) {
	l := uint32(len(*x))
	w.VarUint32(&l)
	_, _ = w.w.Write(*x)
}

func (w *Writer) String(x *string) {
	l := uint32(len(*x))
	w.VarUint32(&l)
	_, _ = w.w.Write(str2Bytes(*x))
}

func (w *Writer) Bool(x *bool) {
	_ = w.w.WriteByte(*(*byte)(unsafe.Pointer(x)))
}

func (w *Writer) Uint8(x *uint8) {
	_ = w.w.WriteByte(*x)
}

func (w *Writer) Int8(x *int8) {
	_ = w.w.WriteByte(byte(*x) & 0xff)
}

func (w *Writer) VarUint64(x *uint64) {
	u := *x
	for u >= 0x80 {
		_ = w.w.WriteByte(byte(u) | 0x80)
		u >>= 7
	}
	_ = w.w.WriteByte(byte(u))
}
