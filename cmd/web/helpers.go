package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Print(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// 根据页面名称（如'home.html'）从缓存中检索适当的模板集合。
	// 如果缓存中不存在具有所提供名称的条目，则创建一个新错误并调用
	// 之前制作的serverError()辅助方法并返回。
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("模板 %s 不存在", page)
		app.serverError(w, err)
		return
	}

	// 初始化新缓冲区
	buf := new(bytes.Buffer)

	// 将模板写入缓冲区，而不是直接写入http.ResponseWriter
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 返回提供的HTTP状态码
	w.WriteHeader(status)

	// 将缓冲区的内容写入http.ResponseWriter
	buf.WriteTo(w)
}
