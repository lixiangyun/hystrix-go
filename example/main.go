package main

import (
	"fmt"
	"time"

	"github.com/lixiangyun/hystrix-go/hystrix"
)

type object struct {
	value int
}

func (o *object) Execute(input interface{}) hystrix.RESULT_TYPE {

	return hystrix.SUCCESS
}

func (o *object) FallBack(input interface{}) {

}

func main() {

	var obj object

	obj.value = 1

	err := hystrix.RegisterCmd("object", &obj)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		err = hystrix.ExecuteCmd("object", i)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	err = hystrix.UnRegisterCmd("object")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	return

}
