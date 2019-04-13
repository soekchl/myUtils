package simpleFileSystem

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	. "github.com/soekchl/myUtils"
)

const (
	html = `
<html>
 <head> 
  <title>免费共享文件服务</title> 
  <link rel="icon" href="data:;base64,=">       <!--先禁止请求网站favicon-->
  <meta name="apple-mobile-web-app-capable" content="yes" /> 
  <meta name="apple-mobile-web-app-status-bar-style" content="black" /> 
  <meta name="viewport" content="width=device-width,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" />
 </head> 
 <body> 
  <details open=""> 
   <summary> <font size="5" title="upload file" >上传文件</font></summary>
   <br /> 
   <form action="#" method="post" enctype="multipart/form-data"> 
    <input type="file" name="uploadFile" /> 
    <input type="submit" title="click upload file" value="点击上传"/>
   </form> 
  </details> 
  <details open=""> 
   <summary> <font size="5" title="share file(click download)" >共享文件(可点击下载)</font></summary>
	%v
   <br /> 
  </details>   
 </body>
</html>
	`
)

var (
	uploadPath    string = "./"             // 目录
	uploadMaxSize int64  = 1024 * 1024 * 10 // 10MB
)

func Start(port, path string, limitUploadSizeMB int64) (err error) {
	err = checkPath(path)
	if err != nil {
		return err
	}
	if limitUploadSizeMB < 1 {
		return errors.New("Limit Upload Size Is Too Small")
	}
	uploadMaxSize = limitUploadSizeMB * 1024 * 1024
	Warn("Port=", port, "\tShare Path=", uploadPath, "\tMax Upload Size=", limitUploadSizeMB, " MB")

	http.HandleFunc("/", mainServerHandle)
	Notice("http://", getIp(), port)
	return http.ListenAndServe(port, nil)
}

func StartGo(port, path string) (err error) {
	err = checkPath(path)
	if err != nil {
		return err
	}
	Warn("Port=", port, " Share Path=", uploadPath)

	ec := make(chan error)
	go func() {
		http.HandleFunc("/", mainServerHandle)
		Notice("http://", port)
		ec <- http.ListenAndServe(port, nil)
	}()

	select {
	case err = <-ec:
	case <-time.After(time.Second / 10):
	}
	close(ec)
	return
}

func mainServerHandle(w http.ResponseWriter, r *http.Request) {
	//从请求当中判断方法
	if r.Method == "GET" {
		if len(r.URL.Path) == 1 {
			showMain(w)
		} else {
			downloadFile(w, r)
		}
	} else {
		uploadFile(w, r)
	}
}

func showMain(w http.ResponseWriter) {
	io.WriteString(w, fmt.Sprintf(html, getShareFileHtml(uploadPath)))
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[1:]
	file, err := os.Open(fileName)
	if err != nil {
		if strings.Index(fileName, ".ico") < 0 {
			Error(err)
		}
		return
	}
	io.Copy(w, file)
	defer file.Close()
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// getFileInfo
	file, head, err := r.FormFile("uploadFile")
	if err != nil {
		Error(err)
		return
	}
	// check size
	if head.Size < 1 {
		io.WriteString(w, fmt.Sprintf("<script>alert('Choice File');window.location.href='/'</script>"))
		return
	}
	if head.Size > uploadMaxSize {
		io.WriteString(w, fmt.Sprintf("<script>alert('File Size Is Too Large');window.location.href='/'</script>"))
		return
	}
	defer file.Close()
	// createFile
	fW, err := os.Create(fmt.Sprint(uploadPath, head.Filename))
	if err != nil {
		io.WriteString(w, fmt.Sprintf("<script>alert('File Create Failed');window.location.href='/'</script>"))
		return
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("<script>alert('File Save Failed');window.location.href='/'</script>"))
		return
	}
	io.WriteString(w, fmt.Sprintf("<script>alert('%v upload OK');window.location.href='/'</script>", head.Filename))
}

func getShareFileHtml(filePath string) (result string) {
	dir, err := ioutil.ReadDir(uploadPath)
	if err != nil {
		Error(err)
		return
	}

	str := ""

	for _, v := range dir {
		if !v.IsDir() {
			str = fmt.Sprintf(`%s<a title="click download file" href="%v">%v</a><br>`, str, v.Name(), v.Name())
		}
	}

	return fmt.Sprintf("<pre>\n%v</pre>", str)
}

// 检查文件路径
func checkPath(path string) error {
	if len(path) < 1 {
		return errors.New("Path Is Error")
	}
	_, err := ioutil.ReadDir(uploadPath)
	if err != nil {
		return err
	}

	n := len(path)
	if path[n-1] == '/' {
		uploadPath = path
	} else {
		uploadPath = fmt.Sprint(path, "/")
	}

	return nil
}

func getIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil && !strings.Contains(ipnet.String(), "/16") {
				return ipnet.IP.To4().String()
			}
		}
	}
	return ""

}
