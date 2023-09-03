package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

var wgV1 sync.WaitGroup

func Busi(channel chan int) {
	for t := range channel {
		fmt.Println("go task = ", t, "goroutine count:", runtime.NumGoroutine())
		wgV1.Done()
	}
}

func main() {
	ch := make(chan int)
	chNum := 4
	for i := 0; i < chNum; i++ {
		go Busi(ch)
	}

	taskCnt := math.MaxInt64
	for t := 0; t < taskCnt; t++ {
		sendTask(t, ch)
	}
}

func sendTask(t int, ch chan int) {
	wgV1.Add(1)
	ch <- t
}
