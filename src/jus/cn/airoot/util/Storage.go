package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	. "jus/str"
	"os"
	. "os"
	"path/filepath"
	"strings"
)

type FileAttr struct {
	Time  int64
	Path  File
	Value string
}

func init() {

}

func GetCode(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(">>", err)
	}
	defer f.Close()
	d, _ := ioutil.ReadAll(f)
	return string(d), err
}

func GetBytes(path string) ([]byte, error) {
	f, err := os.Open(path)
	defer f.Close()
	d, _ := ioutil.ReadAll(f)
	return d, err
}

//--------------------------------复制文件夹--------------------------------------

//遍历目录，将文件信息传入通道
func WalkFiles(src string, dest string, unCopy string) {
	filepath.Walk(src,
		func(f string, fi os.FileInfo, err error) error { //遍历目录
			dPath := Substring(f, StringLen(src), -1)

			if dPath == "" {
				return nil
			}
			a, _ := filepath.Abs(f)
			b, _ := filepath.Abs(unCopy)
			if a == b {
				return filepath.SkipDir
			}
			dPath = dest + "/" + dPath

			if fi.IsDir() {
				if CharAt(fi.Name(), 0) != "." { //只复制开头不为点的数据
					os.MkdirAll(dPath, 0777) //建立文件目录
					//WalkFiles(f, dPath, unCopy) //不用递归遍历，filepath.Wal支持逐级遍历
				} else {
					return filepath.SkipDir
				}
			} else {
				fmt.Println("copy", f, dPath)
				CopyFile(dPath, f)
			}

			return nil

		})
}

/**
 * 复制文件
 */
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

/**
 * 复制函数
 * 复制文件或文件夹
 */
func Copy(src string, dest string, unCopy string) {
	WalkFiles(src, dest, unCopy)
}

/**
 * 模糊查询文件路径转到真实路径
 * 返回值为空则代表没有，否则为用
 */
func JUSExist(name string) string {
	if Exist(name) {
		return name
	}

	//获取文件路径
	path := ""
	file := name
	p := LastIndex(name, "/")
	if p != -1 {
		path = Substring(name, 0, p)
		file = Substring(name, p+1, -1)
	}
	file = strings.ToLower(file)
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return ""
	}
	for _, f := range list {
		if strings.ToLower(f.Name()) == file {
			return path + "/" + f.Name()
		}
	}

	return ""
}

/**
 * 判断文件是否存在
 */
func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

//文件做md5校验
func file2Md5(path string) string {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("Storage Open", err)
		return ""
	}

	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		fmt.Println("Storage Copy", err)
		return ""
	}
	return hex.EncodeToString(md5hash.Sum(nil))
}
