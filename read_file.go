package main

import (
	"fmt"
	"github.com/hpcloud/tail"
)

func main() {
	fileName := "./123.log"
	config := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 开始读取的位置
		ReOpen:    true,                                 // 文件移除或被打包，继续跟踪原文件名 tail -F
		MustExist: true,
		Poll:      true, // 支持文件更改时通知
		//RateLimiter: nil,
		Follow: true, // 实时跟踪
	}
	tailfs, err := tail.TailFile(fileName, config)

	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var msg *tail.Line
	var ok bool

	for {
		// ok 用于判断管道是否被关闭，如果关闭就是文件被重置了，需要重新读取新的管道
		msg, ok = <-tailfs.Lines
		if !ok {
			fmt.Println("tail failed close reopen, fileName:", fileName)
			continue
		}
		fmt.Println("text:", msg.Text)

	}
	//fmt.Println("tail end")
}
