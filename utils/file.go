package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// FileWrite 文件如果不存在那么创建,存在那么写入
func FileWrite(filePath string, writeContent []byte) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err == os.ErrNotExist {
		file, err = os.Create(filePath)
		if err != nil {
			return err
		}
	}

	_, err = file.Write(writeContent)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile 单个文件复制到目标
func CopyFile(src, dst string) error {
	var err error
	var srcFile *os.File
	var dstFile *os.File
	var srcInfo os.FileInfo

	if srcFile, err = os.Open(src); err != nil {
		return err
	}
	defer srcFile.Close()

	if dstFile, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir 递归复制整个目录
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFile := path.Join(src, fd.Name())
		dstFile := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = CopyDir(srcFile, dstFile); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = CopyFile(srcFile, dstFile); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

// PathExists 判断路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// PathIsDir 判断路径是否为文件夹
func PathIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// PathIsFile 判断路径是否为文件
func PathIsFile(path string) bool {
	return !PathIsDir(path)
}

// PathMkdirAll 递归创建文件夹
func PathMkdirAll(filePath string) {
	if !PathExists(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}
