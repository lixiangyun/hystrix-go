package hystrix

type Circuit struct {
	b *Bucket

	swi bool
}

func NewCircuit(length int) *Circuit {

	c := new(Circuit)

	c.b = NewBucket(length)

	return c
}

func (c *Circuit) IsOpen() {

	return c.swi
}
