// 多线程文件搜索工具
package main

import (
	"SearchKeyWord/pool"
	"SearchKeyWord/schKey"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getpk(reader *bufio.Reader, dpath *[]string, keyword *string) {
	fmt.Println("输入对应的路径(不同路径用空格分开)：")

	path, err := reader.ReadString('\n') // 整行读取
	if err != nil {
		err.Error()
		return
	}
	*dpath = strings.Fields(path) // Fields去除空元素，保留有效路径

	fmt.Println("请输入要查找的关键词：")
	keywordStr, _ := reader.ReadString('\n')
	*keyword = strings.TrimSpace(keywordStr)
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	// 输入用户选择
	var tempChoose, choose string

	fmt.Println("----选择你要查找关键词的对象---")
	fmt.Println("【文件夹】\t\t【文件】")
	fmt.Print("请输入：")

	tempChoose, _ = reader.ReadString('\n')
	choose = strings.TrimSpace(tempChoose)

	switch choose {
	case "文件夹":
		var dirPath []string
		var keyword string
		getpk(reader, &dirPath, &keyword)

		//创建协程池，以便于分发任务
		mypool := pool.NewPool()
		for _, name := range dirPath {
			dir := name
			task := schKey.SchDir(dir, keyword)
			mypool.Put(task, 1) // 加入协程池
		}

		mypool.Arrange(10) // 安排10个协程搜索
		mypool.Wait()      // 等待协程搜索结束
		fmt.Println("---当前文件夹搜索结束---")

	case "文件":
		var dpath []string
		var keyword string
		getpk(reader, &dpath, &keyword)

		mypool := pool.NewPool()
		for _, name := range dpath {
			fl := name
			task := schKey.SchFile(fl, keyword)
			mypool.Put(task, 1)
		}

		mypool.Arrange(10)
		mypool.Wait()
		fmt.Println("---当前文件搜索结束---")

	default:
		fmt.Print("请输入有效对象（文件夹 或 文件）")
	}

	fmt.Println("进程结束，按回车键退出")
	reader.ReadString('\n')
}
