package key

import (
	"errors"
	"fmt"
	"strings"
)

type Key interface {
	Namespace() string
	Value() string
	fmt.Stringer
}

const MinecraftNamespace string = "minecraft"

func Parse(key string) (Key, error) {
	s := strings.Split(key, ":")
	if len(s) != 2 {
		return nil, errors.New(`count of ":" must be 1`)
	}
	return New(s[0], s[1]), nil
}

func ParseValid(key string) (Key, error) {
	s := strings.Split(key, ":")
	if len(s) != 2 {
		return nil, errors.New(`count of ":" must be 1`)
	}
	return Make(s[0], s[1])
}

// Make returns a new Key where namespace and value was validated.
func Make(namespace, value string) (k Key, err error) {
	if !NamespaceValid(namespace) || !ValueValid(value) {
		return nil, errors.New("invalid namespace or value")
	}
	return &key{namespace, value}, nil
}

func New(namespace, value string) Key {
	return &key{namespace, value}
}

func NamespaceValid(namespace string) bool {
	for _, char := range namespace {
		if !namespaceCharValid(char) {
			return false
		}
	}
	return true
}

func ValueValid(value string) bool {
	for _, char := range value {
		if !valueCharValid(char) {
			return false
		}
	}
	return true
}

func namespaceCharValid(char rune) bool {
	switch char {
	case '_', '-', '.':
		return true
	}
	return (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')
}

func valueCharValid(char rune) bool {
	return namespaceCharValid(char) || char == '/'
}

type key struct {
	namespace, value string
}

func (k *key) Namespace() string {
	return k.namespace
}

func (k *key) Value() string {
	return k.value
}

func (k *key) String() string {
	return fmt.Sprintf("%s:%s", k.namespace, k.value)
}
