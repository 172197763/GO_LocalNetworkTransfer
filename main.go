package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type netInfo struct {
	Url    string `json:"url"`
	Qrcode string `json:"qrcode"`
	IsShow bool   `json:"is_show"`
}

func main() {
	var restConf rest.RestConf
	ip := ""
	conf.MustLoad("./etc/conf.yaml", &restConf)
	s, err := rest.NewServer(restConf)
	if err != nil {
		log.Fatal(err)
		return
	}
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodGet,
		Path:   "/hello/world",
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			netinfo := make([]netInfo, 0)
			// m["des"] = "请选择适合的内网ip进行访问"
			for k, v := range getINet() {
				tmp := netInfo{
					"http://" + v + ":" + strconv.Itoa(restConf.Port) + "/hello/picAdmin",
					"http://" + v + ":" + strconv.Itoa(restConf.Port) + "/upload/qrcode/url_" + strconv.Itoa(k) + ".png",
					true,
				}
				generateQRCode(tmp.Url, "./upload/qrcode/url_"+strconv.Itoa(k)+".png")
				netinfo = append(netinfo, tmp)
			}
			// writer.Header().Add("Content-Type", "text/html")
			// httpx.OkJson(writer, m)
			jsonData := make(map[string]interface{})
			result, err := json.Marshal(netinfo)
			if err != nil {
				panic(err) // 发生错误时进行相应处理
			}
			jsonData["url_arr"] = string(result)
			writer.Write([]byte(view("./view/inet.html", jsonData)))
			// httpx.Ok(writer)
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodGet,
		Path:   "/hello/picAdmin", //pc端管理图片
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			ip = request.Host
			writer.Header().Add("Content-Type", "text/html")
			// writer.Write([]byte("<h1>123</h1>"))
			data := make(map[string]interface{}, 1)
			data["uploadurl"] = "http://" + ip + "/hello/uploadImg"
			data["loadimgurl"] = "http://" + ip + "/hello/getAllPic"
			data["showpicurl"] = "http://" + ip + "/upload/img/"
			data["emptyimgurl"] = "http://" + ip + "/hello/emptyPic"
			data["createtxtqrcodeurl"] = "http://" + ip + "/hello/createTextQrcode"
			writer.Write([]byte(view("./view/admin.html", data)))
			httpx.Ok(writer)
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodGet,
		Path:   "/hello/getAllPic", //获取所有图片名称
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			dir := "./upload/img"
			files, err := os.ReadDir(dir)
			if err != nil {
				fmt.Println(err)
			}
			res := make([]string, 0)
			for _, file := range files {
				res = append(res, file.Name())
			}
			m := make(map[string][]string)
			m["list"] = res
			httpx.OkJson(writer, m)
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodGet,
		Path:   "/hello/showPic", //输出图片
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			params := request.URL.Query()
			imgName := ""
			// 输出每个参数的键值对
			for key, values := range params {
				for _, value := range values {
					fmt.Printf("%s=%s\n", key, value)
					if key == "img" {
						imgName = value
					}
				}
			}

			data, err := os.ReadFile("./upload/img/" + imgName)
			if err != nil {
				fmt.Println(err)
			}
			writer.Header().Add("Content-Type", "image/png")
			writer.Write(data)
			httpx.Ok(writer)
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodPost,
		Path:   "/hello/emptyPic", //清空文件夹
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			// dirPath := "./upload/img/"
			// if err := clearDir(dirPath); err != nil {
			// 	log.Fatal(err)
			// }
			request.ParseForm()
			post_data := request.Form
			content, _ := post_data["pic_name[]"]
			for _, v := range content {
				clearDir("./upload/img/" + v)
			}

			httpx.OkJson(writer, []int{})
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodPost,
		Path:   "/hello/createTextQrcode", //上传文本接口
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			fmt.Println("text内容")
			request.ParseForm()
			post_data := request.Form
			content, _ := post_data["text"]
			file_set_content("./upload/txt/0000.txt", []byte(content[0]))
			url := "http://" + ip + "/upload/txt/0000.txt"
			generateQRCode(url, "./upload/qrcode/0000.png")
			m := make(map[string]interface{})
			m["code"] = 200
			m["url"] = "http://" + ip + "/upload/qrcode/0000.png"
			httpx.OkJson(writer, m)
		},
	})
	s.AddRoute(rest.Route{ // 添加路由
		Method: http.MethodPost,
		Path:   "/hello/uploadImg", //上传文件接口
		Handler: func(writer http.ResponseWriter, request *http.Request) { // 处理函数
			// 获取表单数据
			err := request.ParseMultipartForm(32 << 20) // 设置最大内存限制为32MB
			if err != nil {
				panic("无法解析表单")
			}

			file, handler, err := request.FormFile("file") // "uploaded_file"是HTML表单中文件输入字段的名称
			if err != nil {
				panic("无法获取上传的文件")
			}
			defer file.Close()

			// 将文件保存到本地或进行其他操作
			data, err := io.ReadAll(file)
			file_set_content("./upload/img/"+handler.Filename, data)
			if err != nil {
				panic("无法读取文件内容")
			}

			// 打印文件信息
			fmt.Println("文件名:", handler.Filename)
			fmt.Println("文件类型:", handler.Header["Content-Type"][0])
			fmt.Printf("文件大小: %d bytes\n", len(data))

			// 返回成功消息
			fmt.Fprintf(writer, "已成功接收并处理文件！")
		},
	})
	//这里注册
	dirlevel := []string{":1", ":2", ":3", ":4", ":5", ":6", ":7", ":8"}
	patern := "/upload/"
	dirpath := "./upload/"
	for i := 1; i < len(dirlevel); i++ {
		path := patern + strings.Join(dirlevel[:i], "/")
		//最后生成 /asset
		s.AddRoute(
			rest.Route{
				Method:  http.MethodGet,
				Path:    path,
				Handler: dirhandler(patern, dirpath),
			})

		logx.Infof("register dir  %s  %s", path, dirpath)
	}

	defer s.Stop()
	openUrl("http://127.0.0.1:" + strconv.Itoa(restConf.Port) + "/hello/world")
	s.Start() // 启动服务

}
func isDir(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if info.IsDir() {
		// fmt.Println("It is a directory")
		return true
	} else {
		// fmt.Println("It is a file")
		return false
	}
}

// 清空文件夹
func clearDir(dirPath string) error {
	if !isDir(dirPath) {
		if err := os.Remove(dirPath); err != nil {
			return err
		}
		return nil
	}
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range dir {
		filePath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			// 如果是文件夹，递归删除
			if err := os.RemoveAll(filePath); err != nil {
				return err
			}
		} else {
			// 如果是文件，直接删除
			if err := os.Remove(filePath); err != nil {
				return err
			}
		}
	}
	return nil
}

// dirhandler函数将指定目录映射到指定路径
func dirhandler(patern, filedir string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		handler := http.StripPrefix(patern, http.FileServer(http.Dir(filedir)))
		handler.ServeHTTP(w, req)
	}
}

// 写入文件
func file_set_content(filePath string, content []byte) bool {
	dir, files := filepath.Split(filePath)
	fmt.Printf("Path: %s\n", filePath)
	fmt.Printf("Directory: %s\n", dir)
	fmt.Printf("Filename: %s\n", files)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println(err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return false
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.Write(content)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	pdir, _ := os.Getwd()
	log.Println("保存文件成功" + pdir)
	return true
}

// 读取html内容
func view(path string, data map[string]interface{}) string {
	// 指定要读取的文件路径
	filePath := path // 将此处替换为你自己的文件路径

	// 通过 ioutil.ReadFile 函数读取文件内容并存入字节切片中
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return ""
	}

	// 将字节切片转换成字符串类型
	contentString := string(contentBytes)

	for k, v := range data {
		str, ok := v.(string)
		if ok {
			fmt.Println("Converted string:", str)
			log.Println("日志--注入参数" + "<goval=" + k + ">:" + str)
			contentString = strings.Replace(contentString, "<goval="+k+">", str, -1)
		} else {
			fmt.Println("Conversion failed")
		}
	}
	// 打印文件内容
	return (contentString)
}

// 查找内网ip
func getINet() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	ip_arr := make([]string, 0)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip_arr = append(ip_arr, ipnet.IP.String())
			}
		}
	}
	return ip_arr
}

// 生成qrcode
func generateQRCode(data string, filePath string) error {
	// 调用 qrcode.Encode 方法生成二维码图片
	content, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return err
	}
	file_set_content(filePath, content)
	return nil
}

// 打开网址
func openUrl(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
		break
	case "windows": // Windows
		cmd = exec.Command("start", url)
		break
	default: // 其他操作系统（如Linux）
		cmd = exec.Command("xdg-open", url)
		break
	}

	err := cmd.Run()
	if err != nil {
		fmt.Println("无法打开浏览器:", err)
	}
}
