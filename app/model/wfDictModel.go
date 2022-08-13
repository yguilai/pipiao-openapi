package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WfDictModel = (*customWfDictModel)(nil)

type (
	// WfDictModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWfDictModel.
	WfDictModel interface {
		wfDictModel
	}

	customWfDictModel struct {
		*defaultWfDictModel
	}
)

// NewWfDictModel returns a model for the database table.
func NewWfDictModel(conn sqlx.SqlConn, c cache.CacheConf) WfDictModel {
	return &customWfDictModel{
		defaultWfDictModel: newWfDictModel(conn, c),
	}
}
