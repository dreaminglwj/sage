package resource

import (
	"fmt"
	"io"
	"sync"
)

var closers []io.Closer

// Register 注册资源
func Register(closer io.Closer) {
	closers = append(closers, closer)
	// fmt.Printf("lwj===>register closers: %+v \n", closers)
}

// Release 释放资源
func Release() {
	// fmt.Printf("lwj===>release closers: %+v \n", closers)
	for _, closer := range closers {
		if closer == nil {
			continue
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()
			if err := closer.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		wg.Wait()
	}
}
