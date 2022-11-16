package define

import "database/sql"

const (
	AndWhere  = "AND"
	OrWhere   = "OR"
	LeftJoin  = "LEFT JOIN"
	InnerJoin = "INNER JOIN"
	RightJoin = "RIGHT JOIN"
)

type Db interface {
	// Table 设置表名 Table("user") | Table("user a")
	Table(tableName string) Db
	// Insert 插入数据库
	Insert() (int64, error)
	// Delete 删除数据
	Delete() (int64, error)
	// Update 更新
	Update() (int64, error)
	// QueryRow 单条查询
	QueryRow(rowFn func(row *sql.Row))
	// Count 总数
	Count() int64
	// ColumnInt64 列整数
	ColumnInt64() []int64
	// ColumnString 列字符串
	ColumnString() []string
	// Query 多条查询
	Query(rowsFn func(rows *sql.Rows))
	// SqlExec [插入｜更新｜删除] 返回[插入ID|影响行数]
	SqlExec(s string, args []any) (sql.Result, error)
	// SqlQuery 查询多条数据
	SqlQuery(s string, args []any, rowsFn func(rows *sql.Rows))
	// SqlQueryRow 查询单挑数据
	SqlQueryRow(s string, args []any, rowFn func(row *sql.Row))
	// Field 设置字段名	Field("a", "b", "c")
	Field(field ...string) Db
	// Value 设置值[插入｜更新]	Value("?", "?", "id=90")
	Value(val ...string) Db
	// Args 设置替换字符值[插入|更新|查询｜删除] Args(100, "0,1")
	Args(args ...any) Db
	// Where 条件语句
	Where(opt string, str string, arg ...any) Db
	// AndWhere 条件语句
	AndWhere(str string, arg ...any) Db
	// OrWhere 并且条件语句
	OrWhere(str string, arg ...any) Db
	// Join 连接操作
	Join(opt, name, on string) Db
	// OrderBy 排序
	OrderBy(orderBy ...string) Db
	// GroupBy 分组
	GroupBy(groupBy ...string) Db
	// OffsetLimit 限制行数
	OffsetLimit(offset, limit int64) Db
}
