package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// 里面注入一个sql连接池
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, concent string, expires int) (int, error) {
	// 写入sql语句
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// 此方法返回一个 sql.Result 类型，其中包含有关语句执行时发生情况的一些基本信息
	result, err := m.DB.Exec(stmt, title, concent, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 返回的 ID 是 int64 类型，因此我们在返回之前将其转换为 int 类型。
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	// 这返回一个指向 sql.Row 对象的指针，该对象保存来自数据库的结果。
	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	// 使用 row.Scan() 将 sql.Row 中每个字段的值复制到 Snippet 结构体中对应的字段。
	// 注意，row.Scan 的参数是你想要将数据复制到的位置的“指针”，并且参数的数量必须与语句返回的列数完全相同。
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// 如果查询没有返回任何行，那么 row.Scan() 将返回一个 sql.ErrNoRows 错误
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) { // 注意：原书方法名为 latest（小写），但前面声明和调用是 Latest（大写），应保持大写。参数名原书为 n，通常用 m。
	//  SQL 语句
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// 使用连接池上的 Query() 方法执行我们的 SQL 语句。
	// 返回一个 sql.Rows 结果集，包含我们查询的结果。
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// 我们延迟调用 rows.Close() 以确保在 Latest() 方法返回之前始终正确关闭 sql.Rows 结果集。
	// 这个 defer 语句应该放在检查 Query() 方法返回错误“之后”。
	// 否则，如果 Query() 返回错误，尝试关闭一个 nil 结果集会导致 panic。
	defer rows.Close()

	// 初始化一个空的切片来保存 Snippet 结构体。
	snippets := []*Snippet{}

	// 使用 rows.Next 迭代结果集中的行。
	// 这准备第一行（以及随后的每一行）以便由 rows.Scan() 方法处理。
	// 如果对所有行的迭代完成，结果集会自动关闭自身并释放底层数据库连接。
	for rows.Next() {
		// 创建一个指向新的零值 Snippet 结构体的指针。
		s := &Snippet{}
		// 使用 rows.Scan() 将行中每个字段的值复制到新创建的 Snippet 对象中。
		// 同样，row.Scan() 的参数是将数据复制到的位置的指针，并且参数的数量必须与语句返回的列数完全相同。
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// 将其追加到 snippets 切片中。
		snippets = append(snippets, s)
	}

	// 当 rows.Next() 循环结束时，我们调用 rows.Err() 来检索迭代期间遇到的任何错误。
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
