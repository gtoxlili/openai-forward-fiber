package pool

func AcquireByteArr(size int) []byte {
	return defaultAllocator.Get(size)
}

func ReleaseByteArr(buf []byte) {
	defaultAllocator.Put(buf)
}
