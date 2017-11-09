package hystrix

import (
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

func (s *Stat) Add(s1 Stat) {
	s.success += s1.success
	s.failure += s1.failure
	s.reject += s1.reject
	s.timeout += s1.timeout
}

type Bucket struct {
	lasttime int

	length int
	stat   []Stat
}

func NewBucket(length int) *Bucket {

	b := new(Bucket)
	b.stat = make([]Stat, length)
	b.length = length

	return b
}

func (b *Bucket) Reset() {
	for i, _ := range b.stat {
		b.stat[i].Reset(0)
	}
}

func (b *Bucket) Stat(success, failure, timeout, reject int) {

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

// input value : [0~100]

func (b *Bucket) FailRate(rate int) bool {

	tm := time.Now().Second()
	var temp Stat

	for _, v := range b.stat {
		if v.time+b.length < tm {
			continue
		}
		temp.Add(v)
	}

	failrate := (temp.failure * 100) / (temp.success + temp.failure)

	b.lasttime = tm

	if failrate >= rate {
		return true
	} else {
		return false
	}
}
