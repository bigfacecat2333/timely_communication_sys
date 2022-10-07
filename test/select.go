package main

import "fmt"

func fibonaci(c, quit chan int)  {
	x, y := 1, 1

	for {
		select {
		case c <- x:
			// 如果c可写
			x = y
			y = x + y
		case <- quit:
			// quit可读
			fmt.Println("quit")
			return
		}
	}
}

func main()  {
	c := make(chan int)
	quit := make(chan int)

	go func()  {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}

		quit <- 0
	}()
	
	fibonaci(c, quit)
	
}