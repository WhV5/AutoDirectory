//============================================================
// 描述: V3.0 自动生成程序运行下的目录,文件做了超链接，需要用md生成html
// 作者: Henry
// 日期: 2019-09-09 11:46
// 版权: Henry @Since 2018
//
//============================================================
package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	File *os.File

	//根据不同的操作系统，加不同的斜杠，仅测试windows，其它未测试
	slash string

	//基本路径拆分的数组长度
	flag int

	//文件前空格
	wordspace string

	//存储已经写入的目录结构，防止再写入
	records []string
)

func main() {
	records = make([]string, 20)
	dirpath, err := getCurrentPath()
	HandleError(err)

	File, err = os.Create(dirpath+`\`+"directory.md")
	split := strings.Split(dirpath, slash)
	flag = len(split)
	HandleError(err)
	defer File.Close()

	err = TravDir(dirpath)
	HandleError(err)

}

//获取当前执行程序的目录文件
func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	HandleError(err)
	path, err := filepath.Abs(file)
	HandleError(err)
	i := strings.LastIndex(path, "/")
	slash = "/"
	if i == -1 {
		i = strings.LastIndex(path, "\\")
		slash = "\\"
	}
	if i == -1{
		i = strings.LastIndex(path,`\`)
		slash = `\`
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i]), nil
}

//遍历文件目录
func TravDir(dirpath string) error {
	file, err := os.OpenFile(dirpath, os.O_RDONLY, os.ModePerm)
	HandleError(err)
	defer file.Close()

	infos, err := file.Readdir(-1)
	HandleError(err)
	for _, file := range infos {
		//如果是目录，需要继续遍历
		if file.IsDir() {
			wordspace = ""
			//找到需要遍历的文件夹名称的共性，剔除掉不需要的遍历的文件
			newdirpath := dirpath + slash + file.Name()
			split := strings.Split(newdirpath, slash)
			for i := 0;i <= len(split)-flag;i++{
				j :=flag + i -1
				if len(records) < j{
					_, err := File.WriteString(wordspace+"["+split[j]+"]("+newdirpath+")"+ "  " + "\r\n")
					HandleError(err)

				}else if split[j] != records[i]{
					_, err := File.WriteString(wordspace+"["+split[j]+"]("+newdirpath+")"+ "  " + "\r\n")
					HandleError(err)
					records[i]=split[j]
				}
				wordspace += " "
			}
			HandleError(err)
			err = TravDir(newdirpath)
			HandleError(err)
		}else {
			StorageFile(dirpath,file.Name(), wordspace)
		}
	}
	return nil
}

//储存文件
func StorageFile(dirpath,filename,wordspace string){
	_, err := File.WriteString(wordspace+"["+filename+"]("+dirpath+slash+filename+")"+ "  " + "\r\n")
	if err != nil {
		panic(err)
	}
}

//错误处理
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
