package main

import (
	"fmt"
	"github.com/mayfield-z/ember/internal/pkg/dpdk_driver"
)

func main() {
	dpdk_driver.InitDpdk()
	fmt.Println(111)
}
