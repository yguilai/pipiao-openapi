package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WfEntryModel = (*customWfEntryModel)(nil)

type (
	// WfEntryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWfEntryModel.
	WfEntryModel interface {
		wfEntryModel
	}

	customWfEntryModel struct {
		*defaultWfEntryModel
	}
)

// NewWfEntryModel returns a model for the database table.
func NewWfEntryModel(conn sqlx.SqlConn, c cache.CacheConf) WfEntryModel {
	return &customWfEntryModel{
		defaultWfEntryModel: newWfEntryModel(conn, c),
	}
}
