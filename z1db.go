package gorm

import (
	"gitee.com/myzero1/gotool/z1err"
	"gitee.com/myzero1/gotool/z1mongo"
)

func Z1ToDryRun(db *DB, model interface{}) {
	// callbacks.go
	// for _, f := range p.fns
	if model != nil {
		m, ok := model.(Z1Modeli)
		if ok && m.DBType() == `mongo` {
			db.DryRun = true
		}
	}
}

func Z1ToMongo(db *DB, model interface{}, stmt *Statement) {
	if model != nil {
		m, ok := model.(Z1Modeli)
		if ok && m.DBType() == `mongo` {
			sql := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
			ret, total, err := z1mongo.Sql2Mongo(sql, true)
			z1err.Check(err)
			db.Error = err
			_ = ret
			_ = total

		}
	}
}

type Z1Modeli interface {
	TableName() string
	DBType() string
}

type Z1Model struct {
	*Model
}

func (m *Z1Model) TableName() string {
	return `TableName_placeholder`
}

func (m *Z1Model) DBType() string {
	return `mysql` // mongo
}
