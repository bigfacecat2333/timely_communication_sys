package main

import (
	"fmt"
	"time"
)

func newTask() {
	i := 0
	for {
		i++
		fmt.Printf("go routine i = %d\n", i)
		time.Sleep(1 * time.Second)
	}
	
}

func main() {
	// 创建一个go程（一种协程），去执行函数/匿名函数
	// exit(0) == runtime.Goexit()
	go newTask()

	i := 0

	for {
		i++
		fmt.Printf("main i = %d\n", i)
		time.Sleep(1 * time.Second)
	}
}