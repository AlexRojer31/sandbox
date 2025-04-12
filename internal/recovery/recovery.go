package recovery

import "fmt"

func Recover() {
	if err := recover(); err != nil {
		fmt.Println(err)
		return
	}
}
