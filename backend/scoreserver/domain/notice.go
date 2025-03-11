package domain

import (
	"context"
	"time"
)

type (
	Notice struct {
		slug          string
		title         string
		effectiveFrom time.Time
		markdown      string
	}
	Notices []*Notice
)

type NoticeData struct {
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	EffectiveFrom time.Time `json:"effective_from"`
	Markdown      string    `json:"markdown"`
}

func (n *Notice) Slug() string {
	return n.slug
}

func (n *Notice) Title() string {
	return n.title
}

func (n *Notice) Markdown() string {
	return n.markdown
}

func (n *Notice) EffectiveFrom() time.Time {
	return n.effectiveFrom
}

func (n *Notice) effective(now time.Time) bool {
	return n.effectiveFrom.Before(now)
}

func (n *Notice) Data() *NoticeData {
	return &NoticeData{
		Slug:          n.slug,
		Title:         n.title,
		Markdown:      n.markdown,
		EffectiveFrom: n.effectiveFrom,
	}
}

func (d *NoticeData) parse() (*Notice, error) {
	if d.Slug == "" {
		return nil, NewInvalidArgumentError("slug is required", nil)
	}
	if d.Title == "" {
		return nil, NewInvalidArgumentError("title is required", nil)
	}
	if d.Markdown == "" {
		return nil, NewInvalidArgumentError("markdown is required", nil)
	}
	if d.EffectiveFrom.IsZero() {
		return nil, NewInvalidArgumentError("effective_from is required", nil)
	}

	return &Notice{
		slug:          d.Slug,
		title:         d.Title,
		markdown:      d.Markdown,
		effectiveFrom: d.EffectiveFrom,
	}, nil
}

func ListNotices(ctx context.Context, eff NoticeReader) (Notices, error) {
	data, err := eff.ListNotices(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get descriptive notices")
	}
	notices := make([]*Notice, 0, len(data))
	for _, noticeData := range data {
		notice, err := noticeData.parse()
		if err != nil {
			return nil, err
		}
		notices = append(notices, notice)
	}
	return notices, nil
}

func ListEffectiveNotices(ctx context.Context, now time.Time, eff NoticeReader) (Notices, error) {
	notices, err := ListNotices(ctx, eff)
	if err != nil {
		return nil, err
	}
	effectiveNotices := make([]*Notice, 0, len(notices))
	for _, notice := range notices {
		if notice.effective(now) {
			effectiveNotices = append(effectiveNotices, notice)
		}
	}
	return effectiveNotices, nil
}

func (n Notices) Save(ctx context.Context, eff NoticeWriter) error {
	data := make([]*NoticeData, 0, len(n))
	for _, notice := range n {
		data = append(data, notice.Data())
	}
	if err := eff.SaveNotices(ctx, data); err != nil {
		return WrapAsInternal(err, "failed to save notices")
	}
	return nil
}

type (
	NoticeReader interface {
		ListNotices(ctx context.Context) ([]*NoticeData, error)
	}
	NoticeWriter interface {
		SaveNotices(ctx context.Context, notices []*NoticeData) error
	}
)
