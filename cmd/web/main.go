package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	// 自动将用户输入的命令行参数转换为对应的类型。类似的还有flag.int()
	// 如果不能转换为string 会报错
	addr := flag.String("addr", ":4000", "HTTP network address") // 4000端口号是默认值

	flag.Parse() // 解析命令行参数
	// 现在可以使用 *addr 获取参数值

	// 利用 log.New()创建用户自定义日志
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// 创建 HTTP 路由多路复用器（ServeMux）并配置静态文件服务
	mux := http.NewServeMux()
	// 创建一个用于提供静态文件访问的处理器（http.Handler）
	/**
	http.Dir("./ui/static/")：将本地文件系统中的 ./ui/static/ 目录包装为 http.FileSystem 接口类型，
	用于指定静态文件所在的根目录（例如该目录下可能有 css/style.css、js/app.js 等文件）。
	http.FileServer(...)：基于指定的 http.FileSystem 创建一个文件服务器处理器，
	该处理器会自动处理对静态文件的 HTTP 请求（如 GET 请求），
	并返回对应的文件内容（支持自动索引目录、处理 404 等）。
	**/
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	//将静态文件服务注册到路由多路复用器 mux 上，使得客户端可以通过 /static 前缀的 URL 访问 ./ui/static/ 目录下的静态文件
	// http.StripPrefix("/static", fileServer)：
	// 移除请求 URL 中的 /static 前缀，再将处理后的请求交给 fileServer 处理。
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Addr：赋值为*addr，确保服务器使用与之前相同的网络地址。
	// ErrorLog：赋值为errorLog，使服务器在出现问题时使用自定义日志记录器记录错误。
	// Handler：赋值为mux，保证服务器沿用之前的路由规则。
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  mux,
	}

	InfoLog.Printf("Starting server on %s", *addr) // Information message
	err := srv.ListenAndServe()                    // 更改为addr 的配置，这样可以在命令行中指定端口
	errlog.Fatal(err)
}
