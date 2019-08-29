package runtime2

import (
	"unsafe"
)

// eface represents the interface{}
type eface struct {
	Type *Type
	Data unsafe.Pointer
}

func efaceOf(ep *interface{}) *eface {
	return (*eface)(unsafe.Pointer(ep))
}

// TypeString return the string of given type
func TypeString(e interface{}) string {
	return efaceOf(&e).Type.String()
}

func Hash(e interface{}) uintptr {
	ep := efaceOf(&e)
	return ep.Type.alg.Hash(ep.Data, 0xdeadbeef)
}
