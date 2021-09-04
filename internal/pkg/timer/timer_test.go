package timer

import (
	"fmt"
	"testing"
)

func expiredFunc(expiredTimes int32) {
	fmt.Println(expiredTimes)
}

func cancelFunc() {
	fmt.Println(-1)
}

func TestTimer(t *testing.T) {
}
