package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Concentration-point/SnippetBox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errLog        *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// 自动将用户输入的命令行参数转换为对应的类型。类似的还有flag.int()
	// 如果不能转换为string 会报错
	addr := flag.String("addr", ":4000", "HTTP network address") // 4000端口号是默认值

	dsn := flag.String("dsn", "root:Qyk880329/@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse() // 解析命令行参数
	// 现在可以使用 *addr 获取参数值

	// 利用 log.New()创建用户自定义日志
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errlog.Fatal(err)
	}

	defer db.Close()

	// 初始化新的模板缓存
	templateCache, err := newTemplateCache()
	if err != nil {
		errlog.Fatal(err)
	}

	app := &application{
		errLog:        errlog,
		infoLog:       InfoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// Addr：赋值为*addr，确保服务器使用与之前相同的网络地址。
	// ErrorLog：赋值为errorLog，使服务器在出现问题时使用自定义日志记录器记录错误。
	// Handler：赋值为mux，保证服务器沿用之前的路由规则。
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  app.routes(), // 挂载到routes中去
	}

	InfoLog.Printf("Starting server on %s", *addr) // Information message
	err = srv.ListenAndServe()                     // 更改为addr 的配置，这样可以在命令行中指定端口
	errlog.Fatal(err)
}

// openDB 创建并验证数据库连接池
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// 验证连接是否有效（sql.Open 不会立即建立连接）
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
