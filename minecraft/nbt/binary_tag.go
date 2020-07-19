package nbt

import "fmt"

// Holds a compound binary tag.
//
// Instead of including an entire NBT implementation, it was decided to
// use this "holder" interface instead. This opens the door for platform specific implementations.
type BinaryTagHolder interface {
	fmt.Stringer // Gets the raw string.
}

func NewBinaryTagHolder(value string) BinaryTagHolder {
	return &binaryTagHolder{value: value}
}

type binaryTagHolder struct {
	value string
}

func (b *binaryTagHolder) String() string {
	return b.value
}
