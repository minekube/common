package key

import "fmt"

type Key interface {
	Namespace() string
	Value() string
	fmt.Stringer
}

const MinecraftNamespace string = "minecraft"

func ValidNew(namespace, value string) (k Key, valid bool) {
	if !NamespaceValid(namespace) || !ValueValid(value) {
		return nil, false
	}
	return &key{namespace, value}, true
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
