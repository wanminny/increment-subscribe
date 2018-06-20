package tools

import (
	"fmt"
	"time"
)

//当前时间（人类可读）
func CurrentTime() string  {

	nowTime := fmt.Sprintf("%s",time.Now().Format("2006-01-02 15:04:05"))

	return nowTime
}
