package main

import (
	"html/template"
	"path/filepath"

	"github.com/Concentration-point/SnippetBox/internal/models"
)

// 定义一个templateData类型，作为我们想要传递给HTML模板的
// 任何动态数据的持有结构。
type templateData struct {
	Snippet *models.Snippet
	//包含一个用于保存片段切片的Snippets字段
	Snippets []*models.Snippet
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
		// 将基础模板文件解析为模板集合。
		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

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
