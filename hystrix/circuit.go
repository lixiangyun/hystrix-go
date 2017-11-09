package hystrix

import (
	"log"
	"sync"
	"time"
)

const (
	CIRCUIT_CLOSE = 0
	CIRCUIT_OPEN  = 1
	CIRCUIT_HALF  = 2
)

type STATUS_TYPE int

type Circuit struct {
	b        *Bucket     // 统计
	length   int         // 统计的周期长度（单位秒）
	rate     int         // 失败率%（0-100）
	delay    int         // 断路器开启到半开启的时间（单位秒）
	opentime int         // 开启的时间
	status   STATUS_TYPE // 断路器状态，关闭、开启、半开启
	retrycnt int         // 重试次数，默认容许一次。
	lock     sync.Mutex  // 读写锁
}

func NewCircuit(length int, rate int, delay int) *Circuit {

	if rate > 100 || rate < 0 {

		log.Println("input rate invailed!", rate)
		return nil
	}

	c := new(Circuit)

	c.status = CIRCUIT_CLOSE
	c.length = length
	c.rate = rate
	c.delay = delay
	c.b = NewBucket(length)

	return c
}

func (c *Circuit) IsOpen() bool {

	c.lock.Lock()
	defer c.lock.Unlock()

	tmnow := time.Now().Second()
	bopen := false

	switch c.status {
	case CIRCUIT_CLOSE:
		{
			bopen := c.b.FailRate(c.rate)
			if bopen == true {
				c.b.Reset()
				c.opentime = tmnow
				c.status = CIRCUIT_OPEN

				bopen = true
			}
		}
	case CIRCUIT_OPEN:
		{
			// 判断断路器时间是否超时
			if c.opentime+c.delay > tmnow {
				// 没有达到阈值时间，继续保持开启
				bopen = true
			} else {
				// 达到时间，进入半开启状态，并且返回一次false
				c.status = CIRCUIT_HALF
			}
		}
	case CIRCUIT_HALF:
		{
			// 继续保持开启
			bopen = true
		}
	default:
		{
			c.status = CIRCUIT_CLOSE
			log.Println("circuit status invailed! ", c.status)
		}
	}

	return bopen
}

func (c *Circuit) Success() {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.status == CIRCUIT_HALF {
		c.status = CIRCUIT_CLOSE
	}

	c.b.Stat(1, 0, 0, 0)
}

func (c *Circuit) Failure() {

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.status == CIRCUIT_HALF {
		c.status = CIRCUIT_OPEN
		c.opentime = time.Now().Second()
	} else {
		c.b.Stat(0, 1, 0, 0)
	}
}

func (c *Circuit) Timeout() {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.b.Stat(0, 0, 1, 0)
}

func (c *Circuit) Reject() {

	c.lock.Lock()
	defer c.lock.Unlock()

	c.b.Stat(0, 0, 0, 1)
}
