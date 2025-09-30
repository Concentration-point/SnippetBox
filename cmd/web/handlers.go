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

	app.render(w, http.StatusOK, "home.html", &templateData{
		Snippets: snippets,
	})

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

	app.render(w, http.StatusOK, "view.html", &templateData{
		Snippet: snippet,
	})
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
