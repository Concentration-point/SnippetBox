package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/Concentration-point/SnippetBox/internal/models"
)

// 定义一个templateData类型，作为我们想要传递给HTML模板的
// 任何动态数据的持有结构。
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	//包含一个用于保存片段切片的Snippets字段
	Snippets []*models.Snippet
}

// 创建一个humanDate函数，它返回一个time.Time对象的字符串表示。
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// 初始化一个template.FuncMap对象并将其存储在全局变量中。这本质上是一个
// 字符串键控的映射，充当我们的自定义模板函数名称和函数本身之间的查找表。
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// 初始化一个新映射作为缓存。
	cache := map[string]*template.Template{}

	// 使用filepath.Glob()函数获取匹配模式"./ui/html/pages/*.tmpl"的所有文件路径的切片。
	// 这基本上会给我们一个应用程序'页面'模板的所有文件路径的切片，
	// 像这样：[ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// 逐个循环遍历页面文件路径。
	for _, page := range pages {
		// 从完整文件路径中提取文件名（如'home.tmpl'）
		// 并将其分配给name变量。
		name := filepath.Base(page)
		// template.FuncMap必须在调用ParseFiles()方法之前注册到模板集合中。
		// 这意味着我们必须使用template.New()创建一个空的模板集合，使用Funcs()方法
		// 注册template.FuncMap，然后正常解析文件。
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// // 将基础模板文件解析为模板集合。
		// ts, err := template.ParseFiles("./ui/html/base.html")
		// if err != nil {
		// 	return nil, err
		// }

		// *在此模板集合上*调用 ParseGlob() 以添加任何部分模板。
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		// *在此模板集合上*调用ParseFiles()以添加页面模板。
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// 正常将模板集合添加到映射中...
		cache[name] = ts
	}

	// 返回映射。
	return cache, nil
}

// 创建一个newTemplateData()辅助函数，它返回一个指向templateData结构体的指针，
// 该结构体用当前年份初始化。注意目前在这里没有使用*http.Request参数，
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
