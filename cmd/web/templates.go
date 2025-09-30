package main

import "github.com/Concentration-point/SnippetBox/internal/models"

// 定义一个templateData类型，作为我们想要传递给HTML模板的
// 任何动态数据的持有结构。
type templateData struct {
	Snippet *models.Snippet
}
