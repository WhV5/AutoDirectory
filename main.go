//============================================================
// 描述: V1.4 读取文件下输入的目录,根据Flag做操作
// 作者: Henry
// 日期: 2019-10-08 10:46
// 版权: Henry @Since 2018
//
//============================================================
package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileS struct {
	Instructions string //关键词解释
	Directory    string //文件目录,如果为空，遍历当前目录的前一级目录
	Flag         int    //功能
	ReplaceWorld string //需要去掉或替换的词
}

var (
	File *os.File
	//根据不同的操作系统，加不同的斜杠，仅测试windows，其它未测试
	slash string
	//基本路径拆分的数组长度
	arrInt int
	//文件前空格
	wordspace string
	//存储已经写入的目录结构，防止再写入
	records []string
	Files   FileS
)

func main() {
	records = make([]string, 20)

	//获取当前程序目录
	dirpath, err := getCurrentPath()
	HandleError(err)
	FileRead, err := os.Open(dirpath + slash + "json.map")
	defer FileRead.Close()
	bytes, err := ioutil.ReadAll(FileRead)
	HandleError(err)

	//反序列化文件内容到Files
	err = json.Unmarshal(bytes, &Files)
	HandleError(err)

	if Files.Flag == 0 {
		File, err = os.Create(dirpath + slash + "directory.md")
		split := strings.Split(dirpath, slash)
		arrInt = len(split)
		HandleError(err)
		defer File.Close()
	}
	//处理文件的目录
	err = HandlerDir(Files.Directory, Files.Flag)
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
	if i == -1 {
		i = strings.LastIndex(path, `\`)
		slash = `\`
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0:i]), nil
}

//遍历文件目录
func HandlerDir(dirpath string, flag int) error {
	file, err := os.OpenFile(dirpath, os.O_RDONLY, os.ModePerm)
	HandleError(err)
	defer file.Close()

	infos, err := file.Readdir(-1)
	HandleError(err)

	if flag == 0 { //遍历文件目录,并输出到文件中
		for _, file := range infos {
			//如果是目录，需要继续遍历
			if file.IsDir() {
				wordspace = "\t"
				//找到需要遍历的文件夹名称的共性，剔除掉不需要的遍历的文件
				newdirpath := dirpath + slash + file.Name()
				split := strings.Split(newdirpath, slash)
				for i := 0; i <= len(split)-arrInt; i++ {
					j := arrInt + i - 1
					if len(records) < j {
						_, err := File.WriteString(wordspace + "[" + split[j] + "](" + newdirpath + ")" + "  " + "\r\n")
						HandleError(err)

					} else if split[j] != records[i] {
						_, err := File.WriteString(wordspace + "[" + split[j] + "](" + newdirpath + ")" + "  " + "\r\n")
						HandleError(err)
						records[i] = split[j]
					}
					wordspace += "\t"
				}
				HandleError(err)
				err = HandlerDir(newdirpath, flag)
				HandleError(err)
			} else {
				StorageFile(dirpath, file.Name(), wordspace)
			}
		}
	} else if flag == 1 { //去掉文件中的关键字
		for _, file := range infos {
			if file.IsDir() {
				newdirpath := dirpath + slash + file.Name()
				err = HandlerDir(newdirpath, flag)
				HandleError(err)
			} else {
				str := strings.Replace(file.Name(), Files.ReplaceWorld, "", -1)
				err := os.Rename(dirpath+slash+file.Name(), dirpath+slash+str)
				HandleError(err)
			}
		}
	}
	return nil
}

//储存文件
func StorageFile(dirpath, filename, wordspace string) {
	_, err := File.WriteString(wordspace + "[" + filename + "](" + dirpath + slash + filename + ")" + "  " + "\r\n")
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
