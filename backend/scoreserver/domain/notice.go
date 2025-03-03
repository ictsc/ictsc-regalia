package domain

import (
	"bufio"
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
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

func ListNotices(ctx context.Context, eff NoticeReader) ([]*Notice, error) {
	data, err := eff.ListNotices(ctx)
	if err != nil {
		return nil, WrapAsInternal(err, "failed to get descriptive notices")
	}
	notices := make([]*Notice, 0, len(data))
	for _, d := range data {
		notice := &Notice{
			id:             d.ID,
			path:           d.Path,
			title:          d.Title,
			markdown:       d.Markdown,
			effectiveFrom:  d.EffectiveFrom,
			effectiveUntil: d.EffectiveUntil,
		}
		notices = append(notices, notice)
	}
	return notices, nil
}

type NoticeRawData struct {
	PageID   string
	PagePath string
	Content  string
	Title    *string
}

func (n *Notice) SaveNotice(ctx context.Context, eff NoticeWriter) error {
	// データを保存
	if err := eff.SaveNotice(ctx, n); err != nil {
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

// TODO: growiのロジックをどうするか聞く
func (d *NoticeRawData) parse() (*Notice, error) {
	contentReader := strings.NewReader(d.Content)
	bodyWriter := &strings.Builder{}

	metadata := make(map[string]string)
	if err := parseNotice(contentReader, bodyWriter, metadata); err != nil {
		return nil, err
	}

	var effectiveFrom, effectiveUntil *time.Time
	if metadata["effective_from"] != "" {
		t, err := time.Parse(time.RFC3339, metadata["effective_from"])
		if err == nil {
			effectiveFrom = &t
		}
	}
	if metadata["effective_until"] != "" {
		t, err := time.Parse(time.RFC3339, metadata["effective_until"])
		if err == nil {
			effectiveUntil = &t
		}
	}

	id, err := uuid.NewV4()
	if err != nil {
		return nil, WrapAsInternal(err, "failed to generate ID")
	}

	return &Notice{
		id:             id,
		path:           d.PagePath,
		title:          metadata["title"],
		markdown:       bodyWriter.String(),
		effectiveFrom:  effectiveFrom,
		effectiveUntil: effectiveUntil,
	}, nil
}

// TODO: growiのロジックをどうするか聞く
func parseNotice(r io.Reader, bodyWriter io.Writer, metadata map[string]string) error {
	scanner := bufio.NewScanner(r)
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
			if match := regexp.MustCompile(`^([^:]+):\s*(.+)$`).FindStringSubmatch(line); match != nil {
				metadata[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
			}
			continue
		}

		if _, err := bodyWriter.Write([]byte(line + "\n")); err != nil {
			return WrapAsInternal(err, "failed to write")
		}
	}

	if err := scanner.Err(); err != nil {
		return WrapAsInternal(err, "failed to scan")
	}
	return nil
}

type (
	NoticeReader interface {
		ListNotices(ctx context.Context) ([]*NoticeData, error)
	}
	NoticeWriter interface {
		SaveNotice(ctx context.Context, notice *Notice) error
	}
	NoticeGetter interface {
		GetNoticeByPath(ctx context.Context, pagePath string) (*NoticeRawData, error)
	}
)
