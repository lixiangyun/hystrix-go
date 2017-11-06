package hystrix

import (
	"sync"
	"time"
)

type Stat struct {
	time    int
	success int
	timeout int
	reject  int
	failure int
}

func (s *Stat) Reset(time int) {
	s.time = time
	s.success = 0
	s.timeout = 0
	s.reject = 0
	s.failure = 0
}

type Bucket struct {
	length int
	stat   []Stat

	lock *sync.RWMutex
}

func NewBucket(length int) *Bucket {

	b := new(Bucket)
	b.stat = make([]Stat, length)
	b.lock = new(sync.RWMutex)
	b.length = length

	return b
}

func (b *Bucket) Reset() {
	b.lock.Lock()
	defer b.lock.Unlock()

	for i, _ := range b.stat {
		b.stat[i].Reset(0)
	}
}

func (b *Bucket) Stat(success, failure, timeout, reject int) {
	b.lock.Lock()
	defer b.lock.Unlock()

	tm := time.Now().Second()
	index := tm % b.length

	stat := &b.stat[index]

	if stat.time != tm {
		stat.Reset(tm)
	}

	stat.success += success
	stat.failure += failure
	stat.timeout += timeout
	stat.reject += reject
}



func (b *Bucket)