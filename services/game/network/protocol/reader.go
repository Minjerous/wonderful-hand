package protocol

import (
	"io"
	"math"
	"unsafe"
)

// Reader 将数据读入 buffer 内
type Reader struct {
	r interface {
		io.Reader
		io.ByteReader
	}
}

func NewReader(r interface {
	io.Reader
	io.ByteReader
}) *Reader {
	return &Reader{r: r}
}

var noMagic = make([]byte, len(Magic))

func (r *Reader) Magic() {
	_, _ = r.r.Read(noMagic)
}

// VarUint32 尝试读 5 个字节获取可变 uint
func (r *Reader) VarUint32(x *uint32) {
	var v uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}

		v |= uint32(b&0x7f) << i // 抹去最高位合并
		if b&0x80 == 0 {         // 最高位 置0，说明已经结束
			*x = v
			return
		}
	}
	panic(errVarIntOverflow)
}

// ByteSlice 从 buffer 内读取一个字节切片到 x 中
func (r *Reader) ByteSlice(x *[]byte) {
	var length uint32
	r.VarUint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil && err != io.EOF {
		panic(err)
	}
	*x = data
}

// ByteSlice 从 buffer 内读取一个string到 x 中
func (r *Reader) String(x *string) {
	s := new([]byte)
	r.ByteSlice(s)
	*x = bytes2Str(*s)
}

func (r *Reader) Bool(x *bool) {
	u, err := r.r.ReadByte()
	if err != nil && err != io.EOF {
		panic(err)
	}
	*x = *(*bool)(unsafe.Pointer(&u))
}

func (r *Reader) Uint8(x *uint8) {
	var err error
	*x, err = r.r.ReadByte()
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func (r *Reader) Int8(x *int8) {
	var b uint8
	r.Uint8(&b)
	*x = int8(b)
}

func (r *Reader) VarUint64(x *uint64) {
	var v uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil && err != io.EOF {
			panic(err)
		}

		v |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(errVarIntOverflow)
}
