package main

import (
	"fmt"
	"testing"
)

func TestCellAlign(t *testing.T) {
	fill := "70981400"
	w := 10
	lPad := fmt.Sprintf("!%%-%ds", w/2) // produces "%-55s"  which is pad left string
	rPad := fmt.Sprintf("%%%ds!", w/2)  // produces "%55s"   which is right pad
	fmt.Println(lPad, rPad, fill)

	str := fmt.Sprintf(lPad, fmt.Sprintf(rPad, fill))
	fmt.Println(str)
	str = fmt.Sprintf("!%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(fill))/2, fill))
	fmt.Print(str)
	str = fmt.Sprintf("!%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(fill))/2, fill))
	fmt.Print(str + "!")
}
