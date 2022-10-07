package main

import (
	"fmt"
	"time"
)

func main() {
	// c := make(chan int) // 无缓冲
	c := make(chan int, 2) // 有缓存

	fmt.Println("len =", len(c), "cap =", cap(c))

	go func() {
		defer fmt.Println(("goroutine exit"))
		fmt.Println("goroutine...")
		for i := 0; i < 4; i++ {
			c <- i
			fmt.Println("发送元素", i, "len =", len(c), "cap =", cap(c))
		}
		close(c)
	}()
	time.Sleep(2 * time.Second)
	for i := 0; i < 4; i++ {
		num := <-c
		fmt.Println("num =",num)
	}

	fmt.Println("main go exit")
	// if data, ok := c; ok {
	// 	fmt.Println("err")
	// } else {
	// 	return
	// }
	for data := range c {
		fmt.Println(data)
	}
}