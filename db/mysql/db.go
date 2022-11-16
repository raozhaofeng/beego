package mysql

import (
	"database/sql"
	"fmt"
	"github.com/raozhaofeng/beego/db/define"
	"strings"
)

type Db struct {
	tx             *sql.Tx  //	事物对象
	db             *sql.DB  //	数据对象
	name           string   //	表名
	fields         []string // 字段名
	values         []string // 字段值
	args           []any    // 替换值
	whereStr       string   //	条件语句
	joinStr        string   //	连接语句
	orderByStr     string   // 排序语句
	groupByStr     string   // 分组语句
	offsetLimitStr string   // 限制行数语句
}

// NewDb 新的Db操作
func NewDb(tx *sql.Tx, db *sql.DB) define.Db {
	return &Db{tx: tx, db: db}
}

// Table 设置表名
func (c *Db) Table(name string) define.Db {
	c.name = name
	return c
}

// QueryRow 单条数据查询
func (c *Db) QueryRow(rowFn func(row *sql.Row)) {
	fields := "*"
	if len(c.fields) > 0 {
		fields = strings.Join(c.fields, ",")
	}
	s := "SELECT " + fields + " FROM " + c.name
	s += c.combineSql()
	c.SqlQueryRow(s, c.args, rowFn)
}

// Query 多条数据查询
func (c *Db) Query(rowsFn func(rows *sql.Rows)) {
	fields := "*"
	if len(c.fields) > 0 {
		fields = strings.Join(c.fields, ",")
	}
	s := "SELECT " + fields + " FROM " + c.name
	s += c.combineSql()

	c.SqlQuery(s, c.args, rowsFn)
}

// Count 条件总数
func (c *Db) Count() int64 {
	var nums int64
	countStr := "count(*)"
	if len(c.fields) == 1 && strings.Index(c.fields[0], "count") > -1 {
		countStr = c.fields[0]
	}
	s := "SELECT " + countStr + " FROM " + c.name
	s += c.combineWhere()

	c.SqlQueryRow(s, c.args, func(row *sql.Row) {
		_ = row.Scan(&nums)
	})
	return nums
}

// ColumnInt64 获取列整形
func (c *Db) ColumnInt64() []int64 {
	var data []int64
	c.Query(func(rows *sql.Rows) {
		var tmp int64
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		data = append(data, tmp)
	})
	return data
}

// ColumnString 获取列字符串
func (c *Db) ColumnString() []string {
	var data []string
	c.Query(func(rows *sql.Rows) {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			panic(err)
		}
		data = append(data, tmp)
	})
	return data
}

// SqlQueryRow 单条数据查询
func (c *Db) SqlQueryRow(s string, args []any, rowFn func(row *sql.Row)) {
	row := c.db.QueryRow(s, args...)
	rowFn(row)
}

// SqlQuery 多条数据查询
func (c *Db) SqlQuery(s string, args []any, rowsFn func(rows *sql.Rows)) {
	rows, err := c.db.Query(s, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		rowsFn(rows)
	}
}

// Field 字段名
func (c *Db) Field(field ...string) define.Db {
	c.fields = append(c.fields, field...)
	return c
}

// Value 字段值
func (c *Db) Value(val ...string) define.Db {
	c.values = append(c.values, val...)
	return c
}

// Args 替换值
func (c *Db) Args(arg ...any) define.Db {
	c.args = append(c.args, arg...)
	return c
}

// Where 条件
func (c *Db) Where(opt string, str string, arg ...any) define.Db {
	if c.whereStr == "" {
		opt = ""
	} else {
		opt = " " + opt + " "
	}
	c.whereStr += opt + str
	if arg != nil {
		c.args = append(c.args, arg...)
	}
	return c
}

// AndWhere 条件语句
func (c *Db) AndWhere(str string, arg ...any) define.Db {
	c.Where("AND", str, arg)
	return c
}

// OrWhere 并且条件语句
func (c *Db) OrWhere(str string, arg ...any) define.Db {
	c.Where("OR", str, arg)
	return c
}

// Join 连接操作
func (c *Db) Join(opt, str, on string) define.Db {
	c.joinStr += " " + opt + " " + str + " ON " + on
	return c
}

// LeftJoin 左连接
func (c *Db) LeftJoin(str, on string) define.Db {
	c.Join("LEFT JOIN", str, on)
	return c
}

// InnerJoin 内连接
func (c *Db) InnerJoin(str, on string) define.Db {
	c.Join("INNER JOIN", str, on)
	return c
}

// RightJoin 右连接
func (c *Db) RightJoin(str, on string) define.Db {
	c.Join("RIGHT JOIN", str, on)
	return c
}

// OrderBy 排序操作
func (c *Db) OrderBy(orderBy ...string) define.Db {
	c.orderByStr = " ORDER BY " + strings.Join(orderBy, ",")
	return c
}

// GroupBy 分组操作
func (c *Db) GroupBy(groupBy ...string) define.Db {
	c.groupByStr = " GROUP BY " + strings.Join(groupBy, ",")
	return c
}

// OffsetLimit 限制行数
func (c *Db) OffsetLimit(offset, limit int64) define.Db {
	c.offsetLimitStr = fmt.Sprintf(" LIMIT %d, %d", offset, limit)
	return c
}

// Insert 插入数据
func (c *Db) Insert() (int64, error) {
	s := "INSERT INTO " + c.name + " ("
	s += strings.Join(c.fields, ",")

	//	补齐问号
	if len(c.fields) == 0 {
		for _, _ = range c.fields {
			c.values = append(c.values, "?")
		}
	}

	s += ") VALUES(" + strings.Join(c.values, ",") + ")"
	result, err := c.SqlExec(s, c.args)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// Delete 删除数据
func (c *Db) Delete() (int64, error) {
	s := "DELETE FROM " + c.name
	s += c.combineWhere()
	result, err := c.SqlExec(s, c.args)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

// Update 更新数据
func (c *Db) Update() (int64, error) {
	s := "UPDATE " + c.name + " SET "
	s += strings.Join(c.values, ",")
	s += c.combineWhere()
	result, err := c.SqlExec(s, c.args)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return rows, nil
}

// SqlExec 语句执行
func (c *Db) SqlExec(s string, args []any) (sql.Result, error) {
	if c.tx != nil {
		return c.tx.Exec(s, args...)
	}
	return c.db.Exec(s, args...)
}

// GetTx 获取Tx
func (c *Db) GetTx() *sql.Tx {
	return c.tx
}

// GetDb 获取Db
func (c *Db) GetDb() *sql.DB {
	return c.db
}

// combineWhere 组合where条件
func (c *Db) combineWhere() string {
	if c.whereStr == "" {
		return ""
	}

	return " WHERE " + c.whereStr
}

// combineSql 组合查询语句
func (c *Db) combineSql() string {
	s := ""
	if c.joinStr != "" {
		s += " " + c.joinStr
	}

	if c.whereStr != "" {
		s += " " + c.whereStr
	}

	if c.groupByStr != "" {
		s += " " + c.groupByStr
	}

	if c.orderByStr != "" {
		s += " " + c.orderByStr
	}

	if c.offsetLimitStr != "" {
		s += " " + c.offsetLimitStr
	}
	return s
}
