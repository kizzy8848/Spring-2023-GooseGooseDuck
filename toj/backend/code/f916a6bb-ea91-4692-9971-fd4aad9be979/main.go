package main

import (
	"fmt"
)

// Two Sum
func main() {
	var a int
	fmt.Scanln(&a)

    sum :=0

    for i:=0;i<=a;i++{
        sum += i
    }

	fmt.Println(sum)
}