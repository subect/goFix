package main

import "fmt"

type Animal interface {
	Eat(string) string
}

type Cat struct {
}

func (cat *Cat) Eat(food string) string {
	fmt.Println("cat eat", food)
	return "cat eat " + food
}

func Live() Animal {
	var animal *Cat
	return animal
}

func main() {
	if Live() == nil {
		fmt.Println("nil")
	} else {
		fmt.Println("not nil")
	}
}
