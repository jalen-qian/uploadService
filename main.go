package main

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	Host   = "http://cdn1.jalen-qian.com/"              //域名
	AK     = "D0Kj0t7WZpzAFzcutiuaa16mvy-Yn8CZZFHnF9om" //公钥
	SK     = "7_JvtkPWhKmkMZriSD6JyXlOh_mfiVE-f7qYCwcM" //私钥
	BUCKET = "jalenqian-waibu"                          //空间
	DIR    = "Jalen/"                                   //内容空间上的文件夹
)

var wg sync.WaitGroup

type Result struct {
	storage.PutRet
	err error
	url string
}

func main() {
	//如果参数中没有传文件名，则直接报错
	if len(os.Args) < 2 {
		fmt.Println("Upload Failed:")
		fmt.Println("Please Select File")
		return
	}
	//定义管道，用来接收文件上传的信息
	resChan := make(chan *Result, len(os.Args)-1)
	//开启多个goroutine来上传文件
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			//网络图片，启动转接模块
			client := &http.Client{}
			resp, err := client.Head(arg)
			if err != nil {
				fmt.Printf("Get Head failed,err:%v\n", err)
				return
			}
			if resp != nil {
				fmt.Printf("%#v\n", resp)
			}
			return
		} else {
			wg.Add(1)
			//本地图片，启动上传模块
			f := &FileMsg{FullFileName: arg}
			f.InitFileMsg().GenSecFileName()
			//启动一个协程去上传
			go upload(f, resChan)
		}

	}
	wg.Wait()
	//所有的协程处理都结束，定义一个Result切片来接收管道中的数据
	results := make([]*Result, 0, len(os.Args)-1)
	for i := 0; i < len(os.Args)-1; i++ {
		result := <-resChan
		//有一个失败，都算全部失败
		if result.err != nil {
			//如果失败，打印失败信息
			fmt.Println("Upload Failed:")
			fmt.Println(result.err)
			return
		} else {
			results = append(results, result)
		}
	}
	//成功，打印成功信息
	fmt.Println("Upload Succeed:")
	for _, r := range results {
		fmt.Println(r.url)
	}

}

/**
上传接口
*/
func upload(fileMsg *FileMsg, resChan chan *Result) {
	defer wg.Done()
	key := DIR + fileMsg.SecFileName
	putPolicy := storage.PutPolicy{
		Scope: BUCKET,
	}
	mac := qbox.NewMac(AK, SK)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	//putExtra := storage.PutExtra{
	//	Params: map[string]string{
	//		"x:name": "github logo",
	//	},
	//}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, fileMsg.FullFileName, nil)
	result := &Result{}
	if err != nil {
		result.err = err
		resChan <- result
		return
	}
	result.PutRet = ret
	result.url = Host + ret.Key
	resChan <- result
}


