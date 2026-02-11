// Package schKey 实现一个目录中的文件关键词提取
package schKey

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"strings"
	"sync"
)

var pLock sync.Mutex //用来保证输出单一化的锁

// 扫描一个文件中关键词的函数 （别看了，包自己用的，别和他抢😎）
func searchFile(flpath string, key string) { //所写的协程池只可以支持func()类型的任务
	//定义检索功能
	if len(key) == 0 {
		fmt.Println("请输入有效关键词")
		return
	}

	file, err := os.OpenFile(flpath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("文件打开失败，错误原因：%v\n", err)
		return
	}
	defer file.Close() // 细节操作 defer放在error后面防止panic偷背身😎

	var result []string // 存储输出信息的字符串切片
	once := sync.Once{}

	scanner := bufio.NewScanner(file)
	line := 0      //行数
	judge := false //判断是否存在关键词

	for scanner.Scan() { //扫描整个文本
		line++
		text := scanner.Text()
		if strings.Contains(text, key) {
			once.Do(func() { judge = true }) //执行一次初始化命令
			s := fmt.Sprintf("搜索到关键词【%s】\n出现在文件路径为 %s 的第 [%d] 行\n", key, flpath, line)
			result = append(result, s) // 存储输出信息，延迟输出
		}
	}

	//上锁，保证文件的统一输出
	pLock.Lock()
	for _, str := range result {
		fmt.Println(str)
	}
	pLock.Unlock()

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取失败：%v\n", err)
	} else {
		if !judge {
			fmt.Printf("文件路径为 %s 的文件未检索到关键词\n\n\n", flpath)
		}
	}
}

// SchDir 搜索目录（对外）
func SchDir(dirpath string, key string) func() {
	return func() {
		dir, err := os.ReadDir(dirpath) // 打开目录，获得目录中内容
		if err != nil {                 // 经典的报错检验
			fmt.Println("目录打开失败，错误原因为：", err)
			return
		}

		if len(dir) == 0 { //判断目录是否为空
			fmt.Println("目录为空")
			return
		}

		//建立文件等待组
		waiter := sync.WaitGroup{}
		//建立子目录等待组
		waiterSdir := sync.WaitGroup{}

		for _, name := range dir {
			if name.IsDir() { // 如果得到的是子目录
				waiterSdir.Add(1)

				go func(fname string) {
					defer waiterSdir.Done()
					dpath := filepath.Join(dirpath, fname) //拼接文件路径
					SchDir(dpath, key)()                   //递归搜索子目录
				}(name.Name()) // 参数传入防止闭包陷阱

			} else { // 如果是文件，则执行下一步搜索关键词操作
				waiter.Add(1)

				// 最讨厌闭包陷阱了😡😡
				go func(filename string) {
					defer waiter.Done()
					fpath := filepath.Join(dirpath, filename)
					searchFile(fpath, key)
				}(name.Name())
			}
			waiterSdir.Wait() // 等待递归遍历完成
		}
		waiter.Wait() // 等待所有文件全部扫描完毕
	}
}

// SchFile 搜索文件（对外）
func SchFile(flpath string, key string) func() {
	return func() {
		searchFile(flpath, key)
	}
}
