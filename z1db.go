package gorm

import (
	"encoding/json"
	"strings"

	"gitee.com/myzero1/gotool/z1mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func Z1ToDryRun(db *DB, modelIsMongo bool) {
	// callbacks.go
	// for _, f := range p.fns

	if modelIsMongo {
		db.DryRun = true
	}
}

func Z1ToMongo(db *DB, model interface{}, stmt *Statement, modelIsMongo bool) {
	if modelIsMongo {
		sql := db.Dialector.Explain(stmt.SQL.String(), stmt.Vars...)
		isCount := true
		isMany := true

		// log.Println(`------sql--1--`, sql)

		if strings.HasPrefix(sql, `SELECT `) {
			if !strings.Contains(sql, `count(`) {
				isCount = false
				if !strings.Contains(sql, ` LIMIT `) {
					sql = sql + ` LIMIT 1`
				}
			}
		}

		// log.Println(`------sql--2--`, sql)

		ret, total, action, err := z1mongo.Sql2Mongo(sql, isCount)

		if err != nil {
			db.Error = err
			return
		}

		db.RowsAffected = total

		if action == `select` {
			if isMany {
				db.RowsAffected = int64(len(ret))
				b, err := bson.Marshal(ret)
				if err != nil {
					db.Error = err
					return
				}
				err = bson.Unmarshal(b, model)
				if err != nil {
					db.Error = err
					return
				}
			} else {
				if len(ret) > 0 {
					db.RowsAffected = 1
					b, err := bson.Marshal(ret[0])
					if err != nil {
						db.Error = err
						return
					}
					err = bson.Unmarshal(b, model)
					if err != nil {
						db.Error = err
						return
					}
				}
			}
		}

		db.DryRun = false
	}
}

func Z1ParsingModel(model interface{}) (isMongo bool) {
	if model != nil {
		m, ok := model.(Z1Modeli)
		if ok && m.DBType() == `mongo` {
			isMongo = true
		}
	}

	return
}

func Z1ParsingModelOld(model interface{}) (isMongo, isSlice bool) {
	if model != nil {
		b, err := json.Marshal(model)

		if err != nil {
			return
		}

		bStr := string(b)

		isSlice = strings.HasPrefix(bStr, `[{"`)

		isMongo = strings.Contains(bStr, `"_id":"000000000000000000000000"`)
	}

	return
}

type Z1Model struct {
	Model
	ID        int64              `gorm:"column:id;primarykey" json:"id" bson:"id,omitempty"`                             // sonyflake machineid max 65536 2^16
	CreatedAt int64              `gorm:"column:created_at;autoCreateTime" json:"created_at" bson:"created_at,omitempty"` // 创建时间戳
	UpdatedAt int64              `gorm:"column:updated_at;autoUpdateTime" json:"updated_at" bson:"updated_at,omitempty"` // 更新时间戳
	ID_       primitive.ObjectID `gorm:"-:all" json:"_id" bson:"_id,omitempty"`                                          // for mongodb _id 这个字段是标识，是否使用MongoDB的
}

type Z1Modeli interface {
	DBType() string
}

func (m *Z1Model) DBType() string {
	return `mysql` // for mongo must return mongo
}
