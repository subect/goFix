package handler

import "fmt"

func ConsumeFunc(param []byte) error {
	fmt.Println("收到请求")
	fmt.Println(string(param))
	return nil
}
