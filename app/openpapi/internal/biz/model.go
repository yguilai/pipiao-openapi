package biz

type I18nEntry struct {
    Uk I18nItem `json:"uk"`
    Zh I18nItem `json:"zh"`
}

type I18nItem struct {
    Name        string `json:"name"`
    Description string `json:"description"`
}
