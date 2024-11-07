package requests

import (
	"fmt"
	"time"
)

func GetCode() string {
	return fmt.Sprintf("%06d", int(time.Now().Nanosecond())%1000000)
}

