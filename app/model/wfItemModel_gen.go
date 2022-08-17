// Code generated by goctl. DO NOT EDIT!

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	wfItemFieldNames          = builder.RawFieldNames(&WfItem{})
	wfItemRows                = strings.Join(wfItemFieldNames, ",")
	wfItemRowsExpectAutoSet   = strings.Join(stringx.Remove(wfItemFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	wfItemRowsWithPlaceHolder = strings.Join(stringx.Remove(wfItemFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheWfItemIdPrefix       = "cache:wfItem:id:"
	cacheWfItemKeyLangPrefix  = "cache:wfItem:key:lang:"
	cacheWfItemNameLangPrefix = "cache:wfItem:name:lang:"
)

type (
	wfItemModel interface {
		Insert(ctx context.Context, data *WfItem) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*WfItem, error)
		FindOneByKeyLang(ctx context.Context, key string, lang string) (*WfItem, error)
		FindOneByNameLang(ctx context.Context, name string, lang string) (*WfItem, error)
		Update(ctx context.Context, data *WfItem) error
		Delete(ctx context.Context, id int64) error
	}

	defaultWfItemModel struct {
		sqlc.CachedConn
		table string
	}

	WfItem struct {
		Id          int64     `db:"id"`          // 自增id
		Key         string    `db:"key"`         // 词条键
		Lang        string    `db:"lang"`        // 语言缩写
		Name        string    `db:"name"`        // 词条名称
		Description string    `db:"description"` // 词条说明
		CreateTime  time.Time `db:"create_time"` // 创建时间
		UpdateTime  time.Time `db:"update_time"` // 更新时间
	}
)

func newWfItemModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultWfItemModel {
	return &defaultWfItemModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`wf_item`",
	}
}

func (m *defaultWfItemModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	wfItemIdKey := fmt.Sprintf("%s%v", cacheWfItemIdPrefix, id)
	wfItemKeyLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemKeyLangPrefix, data.Key, data.Lang)
	wfItemNameLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemNameLangPrefix, data.Name, data.Lang)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, wfItemIdKey, wfItemKeyLangKey, wfItemNameLangKey)
	return err
}

func (m *defaultWfItemModel) FindOne(ctx context.Context, id int64) (*WfItem, error) {
	wfItemIdKey := fmt.Sprintf("%s%v", cacheWfItemIdPrefix, id)
	var resp WfItem
	err := m.QueryRowCtx(ctx, &resp, wfItemIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", wfItemRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWfItemModel) FindOneByKeyLang(ctx context.Context, key string, lang string) (*WfItem, error) {
	wfItemKeyLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemKeyLangPrefix, key, lang)
	var resp WfItem
	err := m.QueryRowIndexCtx(ctx, &resp, wfItemKeyLangKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `key` = ? and `lang` = ? limit 1", wfItemRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, key, lang); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWfItemModel) FindOneByNameLang(ctx context.Context, name string, lang string) (*WfItem, error) {
	wfItemNameLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemNameLangPrefix, name, lang)
	var resp WfItem
	err := m.QueryRowIndexCtx(ctx, &resp, wfItemNameLangKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? and `lang` = ? limit 1", wfItemRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, name, lang); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultWfItemModel) Insert(ctx context.Context, data *WfItem) (sql.Result, error) {
	wfItemIdKey := fmt.Sprintf("%s%v", cacheWfItemIdPrefix, data.Id)
	wfItemKeyLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemKeyLangPrefix, data.Key, data.Lang)
	wfItemNameLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemNameLangPrefix, data.Name, data.Lang)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, wfItemRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.Key, data.Lang, data.Name, data.Description)
	}, wfItemIdKey, wfItemKeyLangKey, wfItemNameLangKey)
	return ret, err
}

func (m *defaultWfItemModel) Update(ctx context.Context, newData *WfItem) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	wfItemIdKey := fmt.Sprintf("%s%v", cacheWfItemIdPrefix, data.Id)
	wfItemKeyLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemKeyLangPrefix, data.Key, data.Lang)
	wfItemNameLangKey := fmt.Sprintf("%s%v:%v", cacheWfItemNameLangPrefix, data.Name, data.Lang)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, wfItemRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.Key, newData.Lang, newData.Name, newData.Description, newData.Id)
	}, wfItemIdKey, wfItemKeyLangKey, wfItemNameLangKey)
	return err
}

func (m *defaultWfItemModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheWfItemIdPrefix, primary)
}

func (m *defaultWfItemModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", wfItemRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultWfItemModel) tableName() string {
	return m.table
}