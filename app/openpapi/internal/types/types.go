// Code generated by goctl. DO NOT EDIT.
package types

type TranslateReq struct {
	Name string `path:"name"`
	Lang string `path:"lang"`
}

type TranslateResp struct {
	Result string `json:"result"`
}
