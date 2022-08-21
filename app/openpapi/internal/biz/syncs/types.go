package syncs

import (
	"context"
	"database/sql"
)

type I18nEntry map[string]I18nItem

type I18nItem struct {
	Name string `json:"name"`
	// 先用空接口吧, 有的是字符串, 有的是数组
	Description interface{} `json:"description"`
}

type SyncService interface {
	NeedFetch(ctx context.Context) (string, string, bool)
	StartUpdate(ctx context.Context, downloadUrl, sha string) error
}

type CommonSyncModel interface {
	FindOld(ctx context.Context, newEntry interface{}) (interface{}, error)
	Add(ctx context.Context, newEntry interface{}) (sql.Result, error)
	Modify(ctx context.Context, oldEntry, newEntry interface{}) error
	NeedModify(oldEntry, newEntry interface{}) bool
}
