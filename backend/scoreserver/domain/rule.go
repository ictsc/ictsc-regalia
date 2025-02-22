package domain

import (
	"context"
)

type Rule struct {
	pagePath string
	markdown string
}

func (r *Rule) PagePath() string {
	return r.pagePath
}

func (r *Rule) Markdown() string {
	return r.markdown
}

func FetchRule(ctx context.Context, eff ProblemContentGetter, pagePath string) (*Rule, error) {
	content, err := eff.GetProblemContentByPath(ctx, pagePath)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to fetch rule")
	}
	return &Rule{
		pagePath: content.PagePath,
		markdown: content.Content,
	}, nil
}

func GetRule(ctx context.Context, eff RuleReader) (*Rule, error) {
	data, err := eff.GetRule(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get rule")
	}
	return data.parse(), nil
}

func (r *Rule) Save(ctx context.Context, eff RuleWriter) error {
	if err := eff.SaveRule(ctx, r.Data()); err != nil {
		return WrapAsInternal(err, "failed to save rule")
	}
	return nil
}

type (
	RuleData struct {
		PagePath string
		Markdown string
	}
	RuleReader interface {
		GetRule(ctx context.Context) (*RuleData, error)
	}
	RuleWriter interface {
		RuleReader
		SaveRule(ctx context.Context, data *RuleData) error
	}
)

func (r *RuleData) parse() *Rule {
	return &Rule{
		pagePath: r.PagePath,
		markdown: r.Markdown,
	}
}

func (r *Rule) Data() *RuleData {
	return &RuleData{
		PagePath: r.pagePath,
		Markdown: r.markdown,
	}
}
