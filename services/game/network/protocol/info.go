package protocol

import "errors"

var (
	// Magic 逸一时，误一世，逸久逸久罢已棂
	Magic = []byte{0x01, 0x01, 0x45, 0x14, 0x19, 0x19, 0x08, 0x10}

	errVarIntOverflow = errors.New("var int over flow")
	errStringTooLong  = errors.New("string(bytes) length overflows a 32-bit integer")
)
