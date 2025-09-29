package main

import (
	"errors"
	"fmt"
	// "html/template"
	"net/http"
	"strconv"

	"github.com/Concentration-point/SnippetBox/internal/models"
)

// Http请求处理函数  加载并渲染html文件
// 使home成为结构体app的方法，这样就可以应用app里的属性
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// 路径校验  避免所有映射都接受
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/pages/home.html",
	// 	"./ui/html/partials/nav.html",
	// }

	// // 将定义的多个模板文件解析为一个模板集合（*template.Template），以便后续渲染。
	// /*template.ParseFiles(files...)：text/template 包（的函数，用于解析指定路径的模板文件。
	//   files... 表示将切片元素展开作为参数传入。
	// */
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// // err = ts.Execute(w, nil)
	// // 将解析后的模板集合中的指定模板（这里是 base）渲染到响应中，返回给客户端。
	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
	// w.Write([]byte("Hello from Snippetbox"))
}
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// stringconvert 将字符串转为int
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 将用户重定向到该片段的相关页面。
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
