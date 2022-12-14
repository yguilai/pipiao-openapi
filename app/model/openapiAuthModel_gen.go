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
	openapiAuthFieldNames          = builder.RawFieldNames(&OpenapiAuth{})
	openapiAuthRows                = strings.Join(openapiAuthFieldNames, ",")
	openapiAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(openapiAuthFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	openapiAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(openapiAuthFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheOpenapiAuthIdPrefix    = "cache:openapiAuth:id:"
	cacheOpenapiAuthAppIdPrefix = "cache:openapiAuth:appId:"
)

type (
	openapiAuthModel interface {
		Insert(ctx context.Context, data *OpenapiAuth) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*OpenapiAuth, error)
		FindOneByAppId(ctx context.Context, appId string) (*OpenapiAuth, error)
		Update(ctx context.Context, data *OpenapiAuth) error
		Delete(ctx context.Context, id int64) error
	}

	defaultOpenapiAuthModel struct {
		sqlc.CachedConn
		table string
	}

	OpenapiAuth struct {
		Id         int64     `db:"id"`          // 自增ID
		AppId      string    `db:"app_id"`      // appId
		AppKey     string    `db:"app_key"`     // appKey
		Status     int64     `db:"status"`      // 状态
		CreateTime time.Time `db:"create_time"` // 创建时间
		UpdateTime time.Time `db:"update_time"` // 更新时间
	}
)

func newOpenapiAuthModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultOpenapiAuthModel {
	return &defaultOpenapiAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`openapi_auth`",
	}
}

func (m *defaultOpenapiAuthModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	openapiAuthAppIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthAppIdPrefix, data.AppId)
	openapiAuthIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, openapiAuthAppIdKey, openapiAuthIdKey)
	return err
}

func (m *defaultOpenapiAuthModel) FindOne(ctx context.Context, id int64) (*OpenapiAuth, error) {
	openapiAuthIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthIdPrefix, id)
	var resp OpenapiAuth
	err := m.QueryRowCtx(ctx, &resp, openapiAuthIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", openapiAuthRows, m.table)
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

func (m *defaultOpenapiAuthModel) FindOneByAppId(ctx context.Context, appId string) (*OpenapiAuth, error) {
	openapiAuthAppIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthAppIdPrefix, appId)
	var resp OpenapiAuth
	err := m.QueryRowIndexCtx(ctx, &resp, openapiAuthAppIdKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `app_id` = ? limit 1", openapiAuthRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, appId); err != nil {
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

func (m *defaultOpenapiAuthModel) Insert(ctx context.Context, data *OpenapiAuth) (sql.Result, error) {
	openapiAuthAppIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthAppIdPrefix, data.AppId)
	openapiAuthIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, openapiAuthRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.AppId, data.AppKey, data.Status)
	}, openapiAuthAppIdKey, openapiAuthIdKey)
	return ret, err
}

func (m *defaultOpenapiAuthModel) Update(ctx context.Context, newData *OpenapiAuth) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	openapiAuthAppIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthAppIdPrefix, data.AppId)
	openapiAuthIdKey := fmt.Sprintf("%s%v", cacheOpenapiAuthIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, openapiAuthRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.AppId, newData.AppKey, newData.Status, newData.Id)
	}, openapiAuthAppIdKey, openapiAuthIdKey)
	return err
}

func (m *defaultOpenapiAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheOpenapiAuthIdPrefix, primary)
}

func (m *defaultOpenapiAuthModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", openapiAuthRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultOpenapiAuthModel) tableName() string {
	return m.table
}
