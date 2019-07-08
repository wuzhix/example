package main

import (
	"demo/cryptology"
	"fmt"
)

func main()  {
	str := "hello world"
	fmt.Printf("Md5Sum %s\n", cryptology.Md5Sum(str))
	fmt.Printf("Md5New %s\n", cryptology.Md5New(str))
}
