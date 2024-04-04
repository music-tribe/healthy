package atomic_bool

import "sync/atomic"

type AtomicBool struct {
	flag uint32
}

func (b *AtomicBool) Set(value bool) {
	var i uint32
	if value {
		i = 1
	}
	atomic.StoreUint32(&b.flag, i)
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadUint32(&b.flag) != 0
}
