package try

import (
	"strings"
	"time"

	"golang.org/x/exp/constraints"
	"xorm.io/builder"

	b "github.com/dreaminglwj/sage/internal/storage/repository/builder"
)

type Ordered interface {
	constraints.Ordered | ~bool
}

// Equal 如果value不为空，则添加 dbField = value 的条件
func Equal[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Eq{dbField: value}
}

// NotEqual 如果value不为空，则添加 dbField != value 的条件
func NotEqual[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Neq{dbField: value}
}

// Gt 如果value不为空，则添加 dbField > value 的条件
func Gt[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Gt{dbField: value}
}

// Gte 如果value不为空，则添加 dbField >= value 的条件
func Gte[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Gte{dbField: value}
}

// Lt 如果value不为空，则添加 dbField < value 的条件
func Lt[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Lt{dbField: value}
}

// Lte 如果value不为空，则添加 dbField <= value 的条件
func Lte[T Ordered](dbField string, value *T) builder.Cond {
	if value == nil {
		return nil
	}
	return builder.Lte{dbField: value}
}

// Like 如果value不为空，则添加 dbField like "%${value}%" 的条件
func Like(dbField, value string) builder.Cond {
	if v := strings.TrimSpace(value); v != "" {
		return b.Like{dbField, v}
	}
	return nil
}

// Likes 如果values不为空，则添加 dbField like "%${value1}%" OR dbField like "%${value2}%" 的条件
func Likes(dbField string, values []string) builder.Cond {
	if len(values) == 0 {
		return nil
	}
	var c []builder.Cond
	for _, v := range values {
		c = append(c, b.Like{dbField, v})
	}
	return builder.Or(c...)
}

// Range 如果min不为空，则添加 dbField >= min 的条件;如果max不为空，则添加 dbField <= max 的条件
func Range[T Ordered](dbField string, min, max *T) builder.Cond {
	if min == nil {
		if max == nil {
			return nil
		} else {
			return builder.Lte{dbField: max}
		}
	} else {
		if max == nil {
			return builder.Gte{dbField: min}
		} else {
			return builder.Between{
				Col:     dbField,
				LessVal: min,
				MoreVal: max,
			}
		}
	}
}

// TimeRange 如果start不为空，则添加 dbField >= min 的条件;如果end不为空，则添加 dbField <= max 的条件
func TimeRange(dbField string, start, end time.Time) builder.Cond {
	if start.IsZero() {
		if end.IsZero() {
			return nil
		} else {
			return builder.Lte{dbField: end}
		}
	} else {
		if end.IsZero() {
			return builder.Gte{dbField: start}
		} else {
			return builder.Between{
				Col:     dbField,
				LessVal: start,
				MoreVal: end,
			}
		}
	}
}

// MultiLike 如果value不为空，则添加 dbField1 like "%${value}%" OR dbField2 like "%${value}%" 的条件
func MultiLike(dbFields []string, value string) builder.Cond {
	if v := strings.TrimSpace(value); v != "" {
		var conds []builder.Cond
		for _, field := range dbFields {
			conds = append(conds, b.Like{field, v})
		}
		return builder.Or(conds...)
	}
	return nil
}

// In 如果values不为空，则添加 dbField IN (values) 的条件
func In[T any](dbField string, values *[]T) builder.Cond {
	if values == nil || len(*values) == 0 {
		return nil
	}
	return builder.In(dbField, *values)
}

// NotIn 如果values不为空，则添加 dbField NOT IN (values) 的条件
func NotIn[T any](dbField string, values []T) builder.Cond {
	if len(values) == 0 {
		return nil
	}
	return builder.NotIn(dbField, values)
}
