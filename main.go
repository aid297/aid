package main

import (
	"fmt"
	"strings"
)

func main() {
	a := strings.HasPrefix("required;in==USER-LEVEL-NORMAL,USER-LEVEL-MANAGER", "in")
	fmt.Printf("%v\n", a)
}
