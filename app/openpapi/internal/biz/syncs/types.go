package syncs

import (
	"context"
	"database/sql"
)

type (
	WfEntry struct {
		UniqueName string `json:"uniqueName"` // wf词条全局唯一名称
		Category   string `json:"category"`   // 词条分类
		Name       string `json:"name"`       // 词条英文名
		Tradable   bool   `json:"tradable"`   // 是否能交易
	}

	I18nEntry map[string]I18nItem

	I18nItem struct {
		Name       string `json:"name"`
		SystemName string `json:"systemName"`
	}

	// SyncService 对
	SyncService interface {
		NeedFetch(ctx context.Context) (string, string, bool)
		StartUpdate(ctx context.Context, downloadUrl, sha string) error
	}

	CommonSync[T any] struct {
		FindOld    func(ctx context.Context, newEntry *T) (*T, error)
		Insert     func(ctx context.Context, newEntry *T) (sql.Result, error)
		Modify     func(ctx context.Context, oldEntry, newEntry *T) error
		NeedModify func(oldEntry, newEntry *T) bool
	}
)
