package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"io/ioutil"
	"path"
	"time"
)

func readFile(filePath string) {
	// 完成读文件的并输出的逻辑

	tailConfig := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 开始读取的位置
		ReOpen:    true,                                 // 文件移除或被打包，继续跟踪原文件名 tail -F
		MustExist: true,
		Poll:      true, // 支持文件更改时通知
		//RateLimiter: nil,
		Follow: true, // 实时跟踪
	}

	tailfs, err := tail.TailFile(filePath, tailConfig)

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
			fmt.Println("tail failed close reopen, fileName:", filePath)
			continue
		}
		fmt.Println("text:", msg.Text)

	}
}

func getFileList(fileDir string) []string {
	fs, _ := ioutil.ReadDir(fileDir)
	fileList := make([]string, 0)

	for _, file := range fs {
		if file.IsDir() {
			//fmt.Println(path+file.Name())
			getFileList(path.Join(fileDir, file.Name()))
		} else {
			if file.Name() == "123.log" {
				fmt.Println(path.Join(fileDir, file.Name()))
				fileList = append(fileList, path.Join(fileDir, file.Name()))
			}
		}
	}
	return fileList
}

func main() {
	// 主流程每秒钟检查目录有没有新的文件产生，有的话就加入启动读取文件任务
	fileName := "D:\\goproject\\read_file"
	//readFile(fileName)
	// 创建一个保存上次读取了文件的切片
	//srcFileList := make([]string, 10)
	// 下次读取文件的切片
	newFileList := make([]string, 0)

	fileList := getFileList(fileName)
	fmt.Println(fileList)
	for _, filePath := range fileList {
		go readFile(filePath)
	}
	// 等一秒
	time.Sleep(1000)

	for {
		newFileList := getFileList(fileName)
		time.Sleep(1000)
	}

	//fmt.Println("tail end")
}
