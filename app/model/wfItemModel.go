package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WfItemModel = (*customWfItemModel)(nil)

type (
	// WfItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWfItemModel.
	WfItemModel interface {
		wfItemModel
	}

	customWfItemModel struct {
		*defaultWfItemModel
	}
)

// NewWfItemModel returns a model for the database table.
func NewWfItemModel(conn sqlx.SqlConn, c cache.CacheConf) WfItemModel {
	return &customWfItemModel{
		defaultWfItemModel: newWfItemModel(conn, c),
	}
}
