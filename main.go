package main

import (
	"fmt"

	"github.com/lzcc1024/goNum2CHS/converter"
)

func main() {

	c := converter.NewConverter()
	s := c.Num2Cap(43314.15)
	fmt.Println(s)

}
