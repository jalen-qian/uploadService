package main

import (
	"github.com/uploadService/mPath"
	"math/rand"
	"strings"
	"time"
)

type FileMsg struct {
	FullFileName       string //全路径 C:\jalen\bin\aa.png
	FileSuffix         string //文件类型 .png
	FileDir            string //文件夹 c:\jalen\bin
	FileNameWithSuffix string //带文件类型的文件名 aa.png
	FileName           string //单纯的文件名 aa
	SecFileName        string //通过算法统一加密后的文件名 比如 2020102732232232.png
}

func (f *FileMsg) InitFileMsg() *FileMsg {
	f.FileNameWithSuffix = mPath.Base(f.FullFileName)                   //获取文件名带后缀
	f.FileSuffix = mPath.Ext(f.FileNameWithSuffix)                      //获取文件后缀
	f.FileName = strings.TrimSuffix(f.FileNameWithSuffix, f.FileSuffix) //获取文件名
	return f
}

func (f *FileMsg) GenSecFileName() {
	if f.FileSuffix == "" {
		return
	}
	//当前系统时间 + 一个随机字符串
	//f.SecFileName = time.Now().Format("20060102150405") + strconv.FormatInt(rand.Int63n(1000), 10) + f.FileSuffix
	f.SecFileName = time.Now().Format("20060102150405") + RandStr(10) + f.FileSuffix
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := make([]byte, 0)
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
