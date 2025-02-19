package domain_test

import (
	"context"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestFetchProblemContent(t *testing.T) {
	t.Parallel()

	eff := problemContentGetterFuncs{
		getProblemContentByIDFunc: func(ctx context.Context, pageID string) (*domain.ProblemContentRawData, error) {
			return nil, errors.New("not implemented")
		},
		getProblemContentByPathFunc: func(ctx context.Context, pagePath string) (*domain.ProblemContentRawData, error) {
			if pagePath != "/test" {
				return nil, domain.NewNotFoundError("page", nil)
			}
			return &domain.ProblemContentRawData{
				PageID:   "test",
				PagePath: "/test",
				Content:  problemExample,
			}, nil
		},
	}

	content, err := domain.FetchProblemContentByPath(t.Context(), eff, "/test")
	if err != nil {
		t.Fatal(err)
	}

	if actual, expected := strings.TrimSpace(content.Body()), strings.TrimSpace(bodyExample); actual != expected {
		t.Errorf("unexpected body: got %v, want %v", actual, expected)
	}
	if actual, expected := strings.TrimSpace(content.Explanation()), strings.TrimSpace(explanationExample); actual != expected {
		t.Errorf("unexpected explanation: got %v, want %v", actual, expected)
	}
}

type (
	getProblemContentByIDFunc   func(ctx context.Context, pageID string) (*domain.ProblemContentRawData, error)
	getProblemContentByPathFunc func(ctx context.Context, pagePath string) (*domain.ProblemContentRawData, error)
	problemContentGetterFuncs   struct {
		getProblemContentByIDFunc
		getProblemContentByPathFunc
	}
)

func (f getProblemContentByIDFunc) GetProblemContentByID(ctx context.Context, pageID string) (*domain.ProblemContentRawData, error) {
	return f(ctx, pageID)
}

func (f getProblemContentByPathFunc) GetProblemContentByPath(ctx context.Context, pagePath string) (*domain.ProblemContentRawData, error) {
	return f(ctx, pagePath)
}

const (
	problemExample = `# -----BEGIN 問題アイディア-----
    
## 発案者
trouble maker

## 概要

なんかうまくいかない

## 詳細

なんかうまくいかないので、なんとかしてほしい

## ネットワーク図
    
## 問題環境

## トラブルの原因

## 解決方法
    

# -----END 問題アイディア-----

---

# -----BEGIN 出題時の問題フォーマット-----

## 問題名

なんかうまくいかない

## 概要

なんかうまくいかないので、なんとかしてほしい

## 前提条件

なし

## 初期状態

うまくいかない

## 終了状態

うまくいく

## 接続情報

| ホスト名 | IPアドレス | ユーザ | パスワード|
| --------- | ----------- | ------ | ------------------ |
| Webサーバ | 192.168.1.1 | user | ictsc2022 |

# -----END 出題時の問題フォーマット-----

---

# -----BEGIN 出題時の問題情報(運営用)-----


## 問題コード

AAA

## 解説

これは Hoge が Piyo なことにより発生する問題です。

## 採点基準

- 1. Moge を満たしていれば 100 点


# -----END 出題時の問題情報(運営用)-----`
	bodyExample = `
## 問題名

なんかうまくいかない

## 概要

なんかうまくいかないので、なんとかしてほしい

## 前提条件

なし

## 初期状態

うまくいかない

## 終了状態

うまくいく

## 接続情報

| ホスト名 | IPアドレス | ユーザ | パスワード|
| --------- | ----------- | ------ | ------------------ |
| Webサーバ | 192.168.1.1 | user | ictsc2022 |`
	explanationExample = `
## 問題コード

AAA

## 解説

これは Hoge が Piyo なことにより発生する問題です。

## 採点基準

- 1. Moge を満たしていれば 100 点`
)
