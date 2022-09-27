package util

import "fmt"

func SafeGo(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("err: ", err)
			}
		}()

		fn()
	}()
}
