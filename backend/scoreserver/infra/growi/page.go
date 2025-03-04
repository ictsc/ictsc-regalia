package growi

import (
	"context"
	"net/url"

	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

var _ domain.ProblemContentGetter = (*Client)(nil)

func (c *Client) GetProblemContentByID(ctx context.Context, pageID string) (*domain.ProblemContentRawData, error) {
	path := &url.URL{Path: "/_api/v3/page"}
	q := path.Query()
	q.Set("pageId", pageID)
	path.RawQuery = q.Encode()

	resp, err := get[growiPageResponse](ctx, c, path)
	if err != nil {
		return nil, err
	}
	return resp.content(), nil
}

func (c *Client) GetProblemContentByPath(ctx context.Context, pagePath string) (*domain.ProblemContentRawData, error) {
	path := &url.URL{Path: "/_api/v3/page"}
	q := path.Query()
	q.Set("path", pagePath)
	path.RawQuery = q.Encode()

	resp, err := get[growiPageResponse](ctx, c, path)
	if err != nil {
		return nil, err
	}
	return resp.content(), nil
}

func (c *Client) GetNoticeByPath(ctx context.Context, pagePath string) (*domain.NoticeRawData, error) {
	path := &url.URL{Path: "/_api/v3/page"}
	q := path.Query()
	q.Set("path", pagePath)
	path.RawQuery = q.Encode()

	resp, err := get[growiPageResponse](ctx, c, path)
	if err != nil {
		return nil, err
	}
	return resp.contentNotice(), nil
}

type (
	growiPageResponse struct {
		Page growiPage `json:"page"`
	}
	growiPage struct {
		ID       string            `json:"id"`
		Path     string            `json:"path"`
		Revision growiPageRevision `json:"revision"`
	}
	growiPageRevision struct {
		Body string `json:"body"`
	}
)

func (r *growiPageResponse) content() *domain.ProblemContentRawData {
	return &domain.ProblemContentRawData{
		PageID:   r.Page.ID,
		PagePath: r.Page.Path,
		Content:  r.Page.Revision.Body,
	}
}

func (r *growiPageResponse) contentNotice() *domain.NoticeRawData {
	return &domain.NoticeRawData{
		PageID:   r.Page.ID,
		PagePath: r.Page.Path,
		Content:  r.Page.Revision.Body,
	}
}
