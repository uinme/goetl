package main

import (
	"etl/liquid"
	"fmt"
)

func main() {
	parserd, err := liquid.Parse("C:Users/uinme/go_workspace/etl/t_goetl_test.yml.liquid")
	if err != nil {
		fmt.Errorf("%w", err)
	}

	fmt.Println(parserd)
}
