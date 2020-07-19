package util

type Codec interface {
	Decode(a interface{}, data []byte) error
}
