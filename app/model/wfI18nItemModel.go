package model

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ WfI18nItemModel = (*customWfI18nItemModel)(nil)

type (
	// WfI18nItemModel is an interface to be customized, add more methods here,
	// and implement the added methods in customWfI18nItemModel.
	WfI18nItemModel interface {
		wfI18nItemModel
		FindOneByNameLang(ctx context.Context, name, lang string) (*WfI18nItem, error)
	}

	customWfI18nItemModel struct {
		*defaultWfI18nItemModel
	}
)

// NewWfI18nItemModel returns a model for the database table.
func NewWfI18nItemModel(conn sqlx.SqlConn, c cache.CacheConf) WfI18nItemModel {
	return &customWfI18nItemModel{
		defaultWfI18nItemModel: newWfI18nItemModel(conn, c),
	}
}

func (m *customWfI18nItemModel) FindOneByNameLang(ctx context.Context, name, lang string) (*WfI18nItem, error) {
	var resp WfI18nItem
	query := fmt.Sprintf("select %s from %s where `name` = ? and `lang` = ? limit 1", wfI18nItemRows, m.table)
	if err := m.CachedConn.QueryRowNoCacheCtx(ctx, &resp, query, name, lang); err != nil {
		return nil, err
	}
	return nil, nil
}
