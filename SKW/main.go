// 多线程文件搜索工具
package main

import (
	"Goland/cmd/goTeam5/SearchKeyWord/pool"
	"Goland/cmd/goTeam5/SearchKeyWord/schKey"
	"fmt"
	"strings"
)

// 获取对象的路径和要查找的关键词
func getpk(dpath *[]string, keyword *string) {
	fmt.Println("输入对应的路径(不同路径用空格分开)：")
	//用bufio.NewReader读取整行
	reader := bufio.NewReader(os.Stdin)
	path, errP := reader.ReadString('\n') // 整行读取
	if errP != nil {
		errP.Error()
		return
	}
	*dpath = strings.Fields(path)     // Fields去除空元素，保留有效路径

	fmt.Println("请输入要查找的关键词：")
	keywordStr, errK := reader.ReadString('\n')
	if errK != nil{
		errK.Error()
		return
	}
	*keyword = strings.TrimSpace(keywordStr)
}

func main() {
	// 输入用户选择
	var choose string
	fmt.Println("----选择你要查找关键词的对象---")
	fmt.Println("【目录】             【文件】")
	fmt.Print("请输入：")
	fmt.Scan(&choose)

	switch choose {
	case "目录":
		var dpath []string
		var keyword string
		getpk(&dpath, &keyword)

		//创建协程池，以便于分发任务
		mypool := pool.NewPool()
		for _, name := range dpath {
			dir := name //可恶的闭包陷阱😡
			task := schKey.SchDir(dir, keyword)
			mypool.Put(task, 1) // 加入协程池
		}

		mypool.Arrange(10) // 安排10个协程搜索
		mypool.Wait()      // 等待协程搜索结束
		fmt.Println("---搜索结束---")

	case "文件":
		var dpath []string
		var keyword string
		getpk(&dpath, &keyword)

		mypool := pool.NewPool()
		for _, name := range dpath {
			fl := name
			task := schKey.SchFile(fl, keyword)
			mypool.Put(task, 1)
		}

		mypool.Arrange(10)
		mypool.Wait()
		fmt.Println("---搜索结束---")

	default:
		fmt.Print("请输入有效对象（目录 或 文件）")
	}
}

