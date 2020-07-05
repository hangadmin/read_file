package main

import (
	"fmt"
	"github.com/hpcloud/tail"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

// 用于保存读到的文件列表内容
var srcFileList []string

const (
	attack = iota + 1
	log
	scan
	crash
)

type RequestData struct {
	eventType int
	scened    string
	grade     string
	sex       string
	address   string
	//string  // 匿名字段，
}

func readFile(filePath, remoteAddr string) {
	// 完成读文件的并输出的逻辑

	tailConfig := tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}, // 开始读取的位置
		ReOpen:    false,                                         // 文件不存在，进程退出
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
		// TODO:处理文件内容
		fmt.Println("text:", msg.Text)
		sendRequest(remoteAddr, attack, "34356", "3445", "攻击消息", "192.168.2.3",
			"192.168.4.7")
	}
}

func getFileList(fileDir string) []string {
	fs, _ := ioutil.ReadDir(fileDir)
	for _, file := range fs {
		if file.IsDir() {
			//fmt.Println(path+file.Name())
			getFileList(path.Join(fileDir, file.Name()))
		} else {
			if file.Name() == "eve.json" {
				//fmt.Println(path.Join(fileDir, file.Name()))
				srcFileList = append(srcFileList, path.Join(fileDir, file.Name()))
			}
		}
	}
	return srcFileList
}

func contractSlice(newFileList, srcFileList []string) []string {
	// 如果新切片中比原切片中多出元素，把多出的值以切片形式返回
	addFileList := make([]string, 0)
out:
	for _, newFile := range newFileList {
		for _, srcFile := range srcFileList {
			if newFile == srcFile {
				continue out
			}
		}
		addFileList = append(addFileList, newFile)
	}
	return addFileList
}

func sendRequest(remoteAddr string, eventType int, sceneId, eventTime, eventMessage, sourceIp, targetIp string) {
	//requestUrl := "http://127.0.0.1:5000/event/push"
	//data := `{"type": attack, "scene_id": "123", "event_time": "", "event_message": "", "source_ip":"", "target_ip": ''}`
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm(remoteAddr, url.Values{
		"type":          {strconv.Itoa(eventType)},
		"scene_id":      {sceneId},
		"event_time":    {eventTime},
		"event_message": {eventMessage},
		"source_ip":     {sourceIp},
		"target_ip":     {targetIp},
	})
	if err != nil {
		fmt.Printf("post data error:%v\n", err)
	} else {
		fmt.Println("post a data successful.")
		respBody, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("response data:%v\n", string(respBody))
	}
}

func main() {
	// 第一个参数指定读取的目录， 第二个参数为YX服务的地址
	fileName := os.Args[1]
	remoteAddr := os.Args[2]
	// 主流程每秒钟检查目录有没有新的文件产生，有的话就加入启动读取文件任务
	//fileName := "E:\\新建文件夹"
	// 创建一个保存上次读取了文件的切片
	// 下次读取文件的切片
	fileList := getFileList(fileName)
	srcFileList = make([]string, 0)
	var newFileList []string
	fmt.Println("原始切片", fileList)
	for _, filePath := range fileList {
		go readFile(filePath, remoteAddr)
	}
	// 等一秒
	time.Sleep(time.Duration(1) * time.Second)

	for {
		newFileList = getFileList(fileName)
		srcFileList = make([]string, 0)
		fmt.Println("最新读到的", newFileList)
		addFileList := contractSlice(newFileList, fileList)
		for _, filePath := range addFileList {
			go readFile(filePath, remoteAddr)
		}
		fileList = newFileList
		fmt.Println("重新赋值原始切片", fileList)
		time.Sleep(5 * time.Second)
	}

	//fmt.Println("tail end")
}
