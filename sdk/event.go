package sdk

import (
	"context"
	"errors"
	"github.com/cellargalaxy/go_common/util"
	"gorm.io/gorm"
	"strings"
	"time"
)

func NewDefaultGormEventHandle() util.GormLogHandle {
	return NewGormEventHandle([]error{gorm.ErrRecordNotFound}, util.DefaultSqlLen)
}

func NewGormEventHandle(ignoreErrs []error, sqlLen int) util.GormLogHandle {
	return GormEventHandle{IgnoreErrs: ignoreErrs, SqlLen: sqlLen}
}

type GormEventHandle struct {
	IgnoreErrs []error
	SqlLen     int
}

func (this GormEventHandle) Handle(ctx context.Context, begin time.Time, sql string, err error) {
	if err == nil {
		return
	}
	ignoreErrs := this.IgnoreErrs
	for i := range ignoreErrs {
		if errors.Is(err, ignoreErrs[i]) {
			return
		}
	}

	if this.SqlLen > 0 && this.SqlLen < len(sql) {
		sql = sql[:this.SqlLen]
	}
	elapsed := time.Since(begin)

	name := strings.Split(sql, " ")[0]
	data := make(map[string]interface{})
	data["sql"] = sql
	data["elapsed"] = elapsed
	data["err"] = err
	AddEvent(ctx, "db", name, 1, data)
}
