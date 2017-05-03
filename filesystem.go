package wke

import (
	"syscall"
	"unsafe"
)

/*
#include "wke.h"

extern void* file_open(const char* path);
extern void file_close(void* handle);
extern size_t file_size(void* handle);
extern int file_read(void* handle, void* buffer, size_t size);
extern int file_seek(void* handle, int offset, int origin);
*/
import "C"

type FileSystem interface {
	Open(path string) uintptr
	Close(handle uintptr)
	Size(handle uintptr) uint32
	Read(handle uintptr, b []byte) int
	Seek(handle uintptr, offset int, origin int) int
}

type DefaultFileSystem struct{}

var fileSystem FileSystem = DefaultFileSystem{}

func (DefaultFileSystem) Open(path string) uintptr {
	handle, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		return uintptr(1<<(unsafe.Sizeof(uintptr(0))/8) - 1)
	}
	return uintptr(handle)
}

func (DefaultFileSystem) Close(handle uintptr) {
	syscall.Close(syscall.Handle(handle))
}

func (DefaultFileSystem) Size(handle uintptr) uint32 {
	var data syscall.ByHandleFileInformation
	if syscall.GetFileInformationByHandle(syscall.Handle(handle), &data) != nil {
		return 0
	}
	if data.FileSizeHigh > 0 {
		panic("too large file")
	}
	return data.FileSizeLow
}

func (DefaultFileSystem) Read(handle uintptr, b []byte) int {
	if n, err := syscall.Read(syscall.Handle(handle), b); err == nil {
		return n
	}
	return -1
}

func (DefaultFileSystem) Seek(handle uintptr, offset int, origin int) int {
	if newoffset, err := syscall.Seek(syscall.Handle(handle), int64(offset), origin); err == nil {
		return int(newoffset)
	}
	return -1
}

//export goFileOpen
func goFileOpen(path *C.char) unsafe.Pointer {
	return unsafe.Pointer(fileSystem.Open(C.GoString(path)))
}

//export goFileClose
func goFileClose(handle unsafe.Pointer) {
	fileSystem.Close(uintptr(handle))
}

//export goFileSize
func goFileSize(handle unsafe.Pointer) C.size_t {
	sz := fileSystem.Size(uintptr(handle))
	return C.size_t(sz)
}

//export goFileRead
func goFileRead(handle unsafe.Pointer, buffer unsafe.Pointer, size C.size_t) int {
	b := (*[1 << 30]byte)(buffer)[:size:size]
	return fileSystem.Read(uintptr(handle), b)
}

//export goFileSeek
func goFileSeek(handle unsafe.Pointer, offset int, origin int) int {
	return fileSystem.Seek(uintptr(handle), offset, origin)
}

func init() {
	C.wkeSetFileSystem(
		C.FILE_OPEN(C.file_open),
		C.FILE_CLOSE(C.file_close),
		C.FILE_SIZE(C.file_size),
		C.FILE_READ(C.file_read),
		C.FILE_SEEK(C.file_seek))
}

func SetFileSystem(fs FileSystem) {
	fileSystem = fs
}
