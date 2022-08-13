package xgithub

import (
    "context"
    "fmt"
    "github.com/google/go-github/github"
    "golang.org/x/oauth2"
)

func NewGithubClient(token string) *github.Client {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    return github.NewClient(tc)
}

func NewSimpleClient() *github.Client {
    return github.NewClient(nil)
}

func TestGetContent() {
    client := NewGithubClient("")
    ctx := context.Background()
    contents, _, _, err := client.Repositories.GetContents(ctx, "WFCD", "warframe-items", "/data/json/i18n.json", nil)
    if err != nil {
        panic(err)
    }
    fmt.Println(contents)
}
