package hystrix

import (
	"fmt"
	"time"
)

func init() {
	go showinfo()
}

func showinfo() {

	fmt.Println("dashboard task start...")

	for {

		time.Sleep(1 * time.Second)
		for name, cir := range circuit {
			fmt.Println("cmd : ", name, " status : ", cir.status)
		}
	}
}
