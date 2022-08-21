package syncs

import (
	"context"
	"database/sql"
)

type WfEntry struct {
	UniqueName string `json:"uniqueName"` // wf词条全局唯一名称
	Category   string `json:"category"`   // 词条分类
	Name       string `json:"name"`       // 词条英文名
	Tradable   bool   `json:"tradable"`   // 是否能交易
}

type I18nEntry map[string]I18nItem

type I18nItem struct {
	Name       string `json:"name"`
	SystemName string `json:"systemName"`
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
