package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
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

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
