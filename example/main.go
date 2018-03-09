package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/tscholl2/colorjson"
)

func main() {
	s := `{"a":[1,2,null],"b":true}`
	fmt.Printf("\nRegular JSON:\n%s\n\nColored JSON:\n", s)
	colorjson.Print(strings.NewReader(s), os.Stdout)
	fmt.Printf("\n\n")
}
