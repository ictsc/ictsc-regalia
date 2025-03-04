package domain

import (
	"bufio"
	"context"
	"io"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"gopkg.in/yaml.v3"
)

type Notice = notice
type notice struct {
	id             uuid.UUID
	path           string
	title          string
	markdown       string
	effectiveFrom  *time.Time
	effectiveUntil *time.Time
}

type NoticeData struct {
	ID             uuid.UUID
	Path           string
	Title          string
	Markdown       string
	EffectiveFrom  *time.Time
	EffectiveUntil *time.Time
}

func (n *Notice) ID() uuid.UUID {
	return n.id
}

func (n *Notice) Path() string {
	return n.path
}

func (n *Notice) Title() string {
	return n.title
}

func (n *Notice) Markdown() string {
	return n.markdown
}

func (n *Notice) EffectiveFrom() *time.Time {
	return n.effectiveFrom
}

func (n *Notice) EffectiveUntil() *time.Time {
	return n.effectiveUntil
}

func (n *Notice) parseToNoticeData() *NoticeData {
	return &NoticeData{
		ID:             n.id,
		Path:           n.path,
		Title:          n.title,
		Markdown:       n.markdown,
		EffectiveFrom:  n.effectiveFrom,
		EffectiveUntil: n.effectiveUntil,
	}
}

func (d *NoticeData) parse() *Notice {
	return &Notice{
		id:             d.ID,
		path:           d.Path,
		title:          d.Title,
		markdown:       d.Markdown,
		effectiveFrom:  d.EffectiveFrom,
		effectiveUntil: d.EffectiveUntil,
	}
}

func ListNotices(ctx context.Context, eff NoticeReader) ([]*Notice, error) {
	data, err := eff.ListNotices(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get descriptive notices")
	}
	notices := make([]*Notice, 0, len(data))
	for _, noticeData := range data {
		notice := &Notice{
			id:             noticeData.ID,
			path:           noticeData.Path,
			title:          noticeData.Title,
			markdown:       noticeData.Markdown,
			effectiveFrom:  noticeData.EffectiveFrom,
			effectiveUntil: noticeData.EffectiveUntil,
		}
		notices = append(notices, notice)
	}
	return notices, nil
}

type NoticeRawData struct {
	PageID   string
	PagePath string
	Content  string
}

func (n *Notice) SaveNotice(ctx context.Context, eff NoticeWriter) error {
	if err := eff.SaveNotice(ctx, n.parseToNoticeData()); err != nil {
		return WrapAsInternal(err, "failed to save notice")
	}

	return nil
}

func FetchNoticeByPath(ctx context.Context, eff NoticeGetter, path string) (*Notice, error) {
	data, err := eff.GetNoticeByPath(ctx, path)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to fetch notice content by path")
	}
	return data.parse()
}

type FrontMatter struct {
	Title          string    `yaml:"title"`
	EffectiveFrom  time.Time `yaml:"effective_from"`
	EffectiveUntil time.Time `yaml:"effective_until"`
}

var (
	ErrMissingFrontMatter    = NewInvalidArgumentError("frontmatter not found", nil)
	ErrMissingTitle          = NewInvalidArgumentError("missing required field: title", nil)
	ErrMissingEffectiveFrom  = NewInvalidArgumentError("missing required field: effective_from", nil)
	ErrMissingEffectiveUntil = NewInvalidArgumentError("missing required field: effective_until", nil)
)

func (d *NoticeRawData) parse() (*Notice, error) {
	contentReader := strings.NewReader(d.Content)
	bodyWriter := &strings.Builder{}

	metadata, err := parseNotice(contentReader, bodyWriter)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to parse metadata")
	}

	if metadata.Title == "" {
		return nil, ErrMissingTitle
	}

	if metadata.EffectiveFrom.IsZero() {
		return nil, ErrMissingEffectiveFrom
	}
	if metadata.EffectiveUntil.IsZero() {
		return nil, ErrMissingEffectiveUntil
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate ID")
	}

	return &Notice{
		id:             id,
		path:           d.PagePath,
		title:          metadata.Title,
		markdown:       bodyWriter.String(),
		effectiveFrom:  &metadata.EffectiveFrom,
		effectiveUntil: &metadata.EffectiveUntil,
	}, nil
}

func parseNotice(r io.Reader, bodyWriter io.Writer) (*FrontMatter, error) {
	scanner := bufio.NewScanner(r)
	var frontMatterContent strings.Builder
	inMetadata := false
	lineno := 0

	for scanner.Scan() {
		lineno++
		line := scanner.Text()

		if strings.HasPrefix(line, "---") {
			inMetadata = !inMetadata
			continue
		}

		if inMetadata {
			frontMatterContent.WriteString(line + "\n")
		} else {
			if _, err := bodyWriter.Write([]byte(line + "\n")); err != nil {
				return nil, WrapAsInternal(err, "failed to write body")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, WrapAsInternal(err, "failed to scan notice content")
	}

	if frontMatterContent.Len() == 0 {
		return nil, ErrMissingFrontMatter
	}

	var frontMatter FrontMatter
	if err := yaml.Unmarshal([]byte(frontMatterContent.String()), &frontMatter); err != nil {
		return nil, WrapAsInternal(err, "failed to parse frontmatter")
	}

	return &frontMatter, nil
}

type (
	NoticeReader interface {
		ListNotices(ctx context.Context) ([]*NoticeData, error)
	}
	NoticeWriter interface {
		SaveNotice(ctx context.Context, notice *NoticeData) error
	}
	NoticeGetter interface {
		GetNoticeByPath(ctx context.Context, pagePath string) (*NoticeRawData, error)
	}
)
