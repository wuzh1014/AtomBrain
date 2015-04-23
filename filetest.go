// 一个简单的目录复制程序：一个独立的 goroutine 遍历目录，主进程负责将数据写入新目录。
package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	RelPath string
	Size    int64
	IsDir   bool
	Handle  *os.File
}

//复制文件数据
func ioCopy(srcHandle *os.File, dstPth string) (err error) {
	dstHandle, err := os.OpenFile(dstPth, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer srcHandle.Close()
	defer dstHandle.Close()
	_, err = io.Copy(dstHandle, srcHandle)
	return err
}

//遍历目录，将文件信息传入通道
func WalkFiles(srcDir, suffix string, c chan<- *FileInfo) {
	suffix = strings.ToUpper(suffix)
	filepath.Walk(srcDir, func(f string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {
			log.Println("[E]", err)
		}
		fileInfo := &FileInfo{}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			if fh, err := os.OpenFile(f, os.O_RDONLY, os.ModePerm); err != nil {
				log.Println("[E]", err)
			} else {
				fileInfo.Handle = fh
				fileInfo.RelPath, _ = filepath.Rel(srcDir, f) //相对路径
				fileInfo.Size = fi.Size()
				fileInfo.IsDir = fi.IsDir()
			}
			c <- fileInfo
		}
		return err
	})
	close(c) //遍历完成，关闭通道

}

//写目标文件
func WriteFiles(dstDir string, c <-chan *FileInfo) {
	if err := os.Chdir(dstDir); err != nil { //切换工作路径
		log.Fatalln("[F]", err)
	}
	for f := range c {
		if fi, err := os.Stat(f.RelPath); os.IsNotExist(err) { //目标不存在
			if f.IsDir {
				if err := os.MkdirAll(f.RelPath, os.ModeDir); err != nil {
					log.Println("[E]", err)
				}
			} else {
				if err := ioCopy(f.Handle, f.RelPath); err != nil {
					log.Println("[E]", err)
				} else {
					log.Println("[I] CP:", f.RelPath)
				}
			}
		} else if !f.IsDir { //目标存在，而且源不是一个目录
			if fi.IsDir() != f.IsDir { //检查文件名被目录名占用冲突
				log.Println("[E]", "filename conflict:", f.RelPath)
			} else if fi.Size() != f.Size { //源和目标的大小不一致时才重写
				if err := ioCopy(f.Handle, f.RelPath); err != nil {
					log.Println("[E]", err)
				} else {
					log.Println("[I] CP:", f.RelPath)
				}
			}
		}
	}
}
func main_1() {
	files_ch := make(chan *FileInfo, 100)
	go WalkFiles("F:\\wait", ".doc", files_ch) //在一个独立的 goroutine 中遍历文件
	WriteFiles("F:\\wait.bak", files_ch)
}
