package define

import "time"

const (
	OptUpdate = 1 //	更新操作
)

// RangeTimeParam 时间范围参数
type RangeTimeParam struct {
	Form string `json:"form"` //	开始时间
	To   string `json:"to"`   //	结束时间
}

// Pagination 分页
type Pagination struct {
	SortBy      string `json:"sortBy"`      //	排序字段
	Descending  bool   `json:"descending"`  //	是否降序排序
	Page        int64  `json:"page"`        //	当前页数
	RowsPerPage int64  `json:"rowsPerPage"` //	每页显示条数
}

type FilterEmpty struct {
	Db  Db  //	模型Db
	opt int //	操作方法
}

func NewFilterEmpty(db Db) *FilterEmpty {
	return &FilterEmpty{
		Db: db,
	}
}

func (c *FilterEmpty) SetUpdateOpt() *FilterEmpty {
	c.opt = OptUpdate
	return c
}

// String 过滤字符串
func (c *FilterEmpty) String(where string, value string) *FilterEmpty {
	if value != "" && value != "%" && value != "%%" {
		c.execute(where, value)
	}
	return c
}

// Int64 过滤整形
func (c *FilterEmpty) Int64(where string, value int64) *FilterEmpty {
	if value != 0 {
		c.execute(where, value)
	}
	return c
}

// Float64 过滤浮点型
func (c *FilterEmpty) Float64(where string, value float64) *FilterEmpty {
	if value != 0 {
		c.execute(where, value)
	}
	return c
}

// DateTime 过滤时间转时间戳
func (c *FilterEmpty) DateTime(where string, value string, location *time.Location) *FilterEmpty {
	if value != "" {
		var valueTime time.Time
		if len(value) == 19 {
			valueTime, _ = time.ParseInLocation("2006/01/02 15:04:05", value, location)
		} else {
			valueTime, _ = time.ParseInLocation("2006/01/02", value, location)
		}

		c.execute(where, valueTime.Unix())
	}
	return c
}

// RangeTime 是否时间范围
func (c *FilterEmpty) RangeTime(where string, rangeTime *RangeTimeParam, location *time.Location) *FilterEmpty {
	if rangeTime != nil {
		startTime, _ := time.ParseInLocation("2006/01/02 15:04:05", rangeTime.Form, location)
		endTime, _ := time.ParseInLocation("2006/01/02 15:04:05", rangeTime.To, location)
		c.Db.AndWhere(where, startTime.Unix(), endTime.Unix())
	}
	return c
}

// Pagination 分页设置
func (c *FilterEmpty) Pagination(p *Pagination) *FilterEmpty {
	if p != nil {
		if p.Page <= 0 {
			p.Page = 1
		}
		if p.RowsPerPage < 1 {
			p.RowsPerPage = 1
		}
		descending := "asc"
		if p.Descending {
			descending = "desc"
		}
		if p.SortBy != "" {
			c.Db.OrderBy(p.SortBy + " " + descending)
		}
		c.Db.OffsetLimit((p.Page-1)*p.RowsPerPage, p.RowsPerPage)
	}
	return c
}

func (c *FilterEmpty) execute(where string, value any) {
	switch c.opt {
	case OptUpdate:
		c.Db.Value(where).Args(value)
	default:
		c.Db.AndWhere(where, value)
	}
}
