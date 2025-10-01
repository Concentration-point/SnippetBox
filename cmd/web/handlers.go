package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Concentration-point/SnippetBox/internal/models"
)

// Http请求处理函数  加载并渲染html文件
// 使home成为结构体app的方法，这样就可以应用app里的属性
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// // 路径校验  避免所有映射都接受
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }
	// httprouter 精确匹配"/" 路径，去掉上面的代码

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 获取包含默认年份数据的templateData结构体
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

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

	// 获取包含默认年份数据的templateData结构体
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)
}

// 添加一个新的 snippetCreate 处理器，目前返回一个占位符响应。
// 我们稍后会更新它以显示一个 HTML 表单。
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

// 将此处理器重命名为 snippetCreatePost。
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// 检查请求方法是否为 POST 现在是多余的，可以移除，
	// 由 httprouter 自动完成的。

	title := "0 snail"
	content := "0 snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// 更新重定向路径以使用新的干净 URL 格式。
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
