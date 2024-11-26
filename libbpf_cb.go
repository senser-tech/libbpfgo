package libbpfgo

import (
	"C"
	"unsafe"
)

// This callback definition needs to be in a different file from where it is declared in C
// Otherwise, multiple definition compilation error will occur

//export perfCallback
func perfCallback(ctx unsafe.Pointer, cpu C.int, data unsafe.Pointer, size C.int) {
	pb := eventChannels.Get(uint(uintptr(ctx))).(*PerfBuffer)
	ch := pb.eventsChan
	switch ch.(type) {
	case chan []byte:
		ch.(chan []byte) <- C.GoBytes(data, size)
	default:
		data_raw := (*[1 << 30]byte)(unsafe.Pointer(data))[:size:size]
		ch.(HandlerType)(data_raw)
	}
}

//export perfLostCallback
func perfLostCallback(ctx unsafe.Pointer, cpu C.int, cnt C.ulonglong) {
	pb := eventChannels.Get(uint(uintptr(ctx))).(*PerfBuffer)
	if pb.lostChan != nil {
		pb.lostChan <- uint64(cnt)
	}
}

type HandlerType = func(dataRaw []byte)

//export ringbufferCallback
func ringbufferCallback(ctx unsafe.Pointer, data unsafe.Pointer, size C.int) C.int {
	ch := eventChannels.Get(uint(uintptr(ctx)))
	switch ch.(type) {
	case chan []byte:
		ch.(chan []byte) <- C.GoBytes(data, size)
		return C.int(0)
	default:
		data_raw := (*[1 << 30]byte)(unsafe.Pointer(data))[:size:size]
		ch.(HandlerType)(data_raw)
	}

	return C.int(0)
}
