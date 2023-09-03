package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func main() {
	taskCnt := 10

	buffer := make(chan bool, 4)
	for i := 0; i < taskCnt; i++ {
		wg.Add(1)
		buffer <- true
		go Yewu(buffer, i)
	}
	//time.Sleep(500 * time.Second)
	wg.Wait()

}

func Yewu(buffer chan bool, num int) {
	fmt.Println("go func", num, "goroutine count:", runtime.NumGoroutine())
	<-buffer
	wg.Done()
}
