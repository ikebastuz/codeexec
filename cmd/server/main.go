package main

import (
	"codeexec/internal/runner"
	"fmt"
)

func main() {
	out, err := runner.Run("javascript", "var a = 1;\nvar b = 2;\nconsole.log(a+b);")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	fmt.Printf("Output: %s\n", out)
}
