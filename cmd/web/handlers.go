package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Http请求处理函数  加载并渲染html文件
func home(w http.ResponseWriter, r *http.Request) {
	// 路径校验  避免所有映射都接受
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/nav.html",
	}

	// 将定义的多个模板文件解析为一个模板集合（*template.Template），以便后续渲染。
	/*template.ParseFiles(files...)：text/template 包（的函数，用于解析指定路径的模板文件。
	  files... 表示将切片元素展开作为参数传入。
	*/
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// err = ts.Execute(w, nil)
	// 将解析后的模板集合中的指定模板（这里是 base）渲染到响应中，返回给客户端。
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	// w.Write([]byte("Hello from Snippetbox"))
}
func snippetView(w http.ResponseWriter, r *http.Request) {
	// stringconvert 将字符串转为int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
