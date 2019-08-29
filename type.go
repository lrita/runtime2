package runtime2

import (
	"reflect"
	"unsafe"
)

// tflag is documented in reflect/type.go.
//
// tflag values must be kept in sync with copies in:
//	cmd/compile/internal/gc/reflect.go
//	cmd/link/internal/ld/decodesym.go
//	reflect/type.go
type tflag uint8
type nameOff int32
type typeOff int32

const (
	tflagUncommon  tflag = 1 << 0
	tflagExtraStar tflag = 1 << 1
	tflagNamed     tflag = 1 << 2

	kindMask = (1 << 5) - 1
)

// TypeAlg is also copied/used in reflect/type.go.
// keep them in sync.
type TypeAlg struct {
	// function for hashing objects of this type
	// (ptr to object, seed) -> hash
	Hash func(unsafe.Pointer, uintptr) uintptr
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	Equal func(unsafe.Pointer, unsafe.Pointer) bool
}

// Needs to be in sync with ../cmd/link/internal/ld/decodesym.go:/^func.commonsize,
// ../cmd/compile/internal/gc/reflect.go:/^func.dcommontype and
// ../reflect/type.go:/^type.rtype.
type Type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *TypeAlg
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

func (t *Type) nameOff(off nameOff) name {
	return resolveNameOff(unsafe.Pointer(t), off)
}

// String return the type string.
func (t *Type) String() string {
	s := t.nameOff(t.str).name()
	if t.tflag&tflagExtraStar != 0 {
		return s[1:]
	}
	return s
}

// Size return the size of this type.
func (t *Type) Size() int {
	return int(t.size)
}

// Kind represents the type kind.
func (t *Type) Kind() reflect.Kind {
	return reflect.Kind(t.kind & kindMask)
}

// name is an encoded type name with optional extra data.
// See reflect/type.go for details.
type name struct {
	bytes *byte
}

func (n name) data(off int) *byte {
	return (*byte)(add(unsafe.Pointer(n.bytes), uintptr(off)))
}

func (n name) nameLen() int {
	return int(uint16(*n.data(1))<<8 | uint16(*n.data(2)))
}

func (n name) name() (s string) {
	if n.bytes == nil {
		return ""
	}
	nl := n.nameLen()
	if nl == 0 {
		return ""
	}
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	hdr.Data = uintptr(unsafe.Pointer(n.data(3)))
	hdr.Len = nl
	return s
}

//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}

//go:linkname resolveNameOff runtime.resolveNameOff
func resolveNameOff(ptrInModule unsafe.Pointer, off nameOff) name
