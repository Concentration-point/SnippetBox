package main

// 中间件包

import (
	"fmt"
	"net/http"
)

func secureHeaders(next http.Handler) http.Handler {
	// 返回值也是一个函数
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置安全请求头逻辑
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		// 调用下一个处理链
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建一个延迟函数（在发生 panic 且 Go 展开堆栈时，该函数将始终运行）。
		defer func() {
			// 使用内置的 recover 函数检查是否发生了 panic
			if err := recover(); err != nil {
				// 在响应上设置 "Connection: close" 头部。
				w.Header().Set("Connection", "close")
				// 调用 app.serverError 辅助方法返回 500 Internal Server 响应。
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
