package monitor

import (
	"fmt"
	"testing"
)

func TestSetResourceCount(t *testing.T) {
	ok := SetResourceCount("Goroutines", int64(1000))
	fmt.Println("result:", ok)
	fmt.Printf("Goroutines: %#v\n", GetResource("Goroutines"))
}
