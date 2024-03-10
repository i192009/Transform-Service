package libzcad

// #cgo LDFLAGS: -L./build/ -lhello
// #include <stdlib.h>
// #include "hello.h"
import "C"
import (
	"time"
	"unsafe"
)

func Hello(file string) {
	time.Sleep(3 * time.Second)
	p := C.CString(file)
	defer C.free(unsafe.Pointer(p))
	C.Hello(p)
}
