package biz

type I18nEntry map[string]I18nItem

type I18nItem struct {
    Name string `json:"name"`
    // 先用空接口吧, 有的是字符串, 有的是数组
    Description interface{} `json:"description"`
}
