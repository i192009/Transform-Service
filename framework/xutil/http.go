package xutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.zixel.cn/go/framework/logger"
)

var (
	log = logger.Get()
)

type RemoteFileInfo struct {
	FileName   string
	FileSize   int64
	ModifyTime time.Time
	CreateTime time.Time
}

func GetRemoteFileinfo(downloadUrl string, info *RemoteFileInfo) error {
	request, err := http.NewRequest("HEAD", downloadUrl, nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("response status code error")
	}

	o, err := url.Parse(downloadUrl)
	if err != nil {
		return err
	}

	info.FileName = filepath.Base(o.Path)
	info.FileSize, err = strconv.ParseInt(response.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}

	info.ModifyTime, err = time.Parse(response.Header.Get("Last-Modified"), time.RFC1123)
	if err != nil {
		return err
	}

	info.CreateTime, err = time.Parse(response.Header.Get("Date"), time.RFC1123)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFile(downloadUrl string, saveTo string) error {

	var info RemoteFileInfo
	if err := GetRemoteFileinfo(downloadUrl, &info); err != nil {
		return err
	}

	resp, err := http.Get(downloadUrl)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(saveTo)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func UploadFile(url string, path string) error {
	// 打开本地文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	client := &http.Client{}
	method := "PUT"
	req, err := http.NewRequest(method, url, file)

	if err != nil {
		log.Errorf("UploadFile NewRequest is error.err:%s", err)
		return err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Host", "re-link.obs.cn-east-3.myhuaweicloud.com")
	req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Content-Type", "application/octet-stream")

	res, err := client.Do(req)
	if err != nil {
		log.Errorf("UploadFile Do is error.err:%s", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(res.Body)

	return nil
}

func DirectInvoke(host string, method string, reqUri string, reqBody io.Reader, header http.Header) (res []byte, rpnCode int, err error) {
	req, err := http.NewRequestWithContext(context.Background(), strings.ToUpper(method),
		fmt.Sprintf("http://%v%v", host, reqUri), reqBody)
	req.Header = header
	if err != nil {
		return
	}
	var rpn *http.Response
	rpn, err = http.DefaultClient.Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("http invoke failed,error=%v", err))
		return
	}
	defer rpn.Body.Close()

	res, err = io.ReadAll(rpn.Body)
	if err != nil {
		err = errors.New(fmt.Sprintf("read http response failed,error=%v", err))
		return
	}

	log.Debugf("http invoke success,req=%v,rpn=%v", *req, *rpn)

	rpnCode = rpn.StatusCode
	return
}
