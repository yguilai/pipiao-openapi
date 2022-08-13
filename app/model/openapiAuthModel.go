package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OpenapiAuthModel = (*customOpenapiAuthModel)(nil)

type (
	// OpenapiAuthModel is an interface to be customized, add more methods here,
	// and implement the added methods in customOpenapiAuthModel.
	OpenapiAuthModel interface {
		openapiAuthModel
	}

	customOpenapiAuthModel struct {
		*defaultOpenapiAuthModel
	}
)

// NewOpenapiAuthModel returns a model for the database table.
func NewOpenapiAuthModel(conn sqlx.SqlConn, c cache.CacheConf) OpenapiAuthModel {
	return &customOpenapiAuthModel{
		defaultOpenapiAuthModel: newOpenapiAuthModel(conn, c),
	}
}
