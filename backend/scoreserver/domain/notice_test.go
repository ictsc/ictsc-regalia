package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ictsc/ictsc-regalia/backend/scoreserver/domain"
)

func TestSyncNotice(t *testing.T) {
	t.Parallel()

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	// NoticeRawDataを作って、parseする
	cases := map[string]struct {
		inNoticeRawData *domain.NoticeRawData
		wantNotice      *domain.Notice
		wantErr         error
	}{
		"ok": {
			inNoticeRawData: &domain.NoticeRawData{
				PageID:   "test",
				PagePath: "/test",
				Content: `---
title: テストお知らせ
draft: false
effective_from: ` + yesterday.Format(time.RFC3339) + `
effective_until: ` + tomorrow.Format(time.RFC3339) + `
---
これはサンプルです。`,
			},

			wantNotice: domain.FixNotice1(t, nil),
		},
		"title not found": {
			inNoticeRawData: &domain.NoticeRawData{
				PageID:   "test",
				PagePath: "/test",
				Content: `---
draft: false
effective_from: ` + yesterday.Format(time.RFC3339) + `
effective_until: ` + tomorrow.Format(time.RFC3339) + `
---
これはサンプルです`,
			},

			wantErr: domain.ErrMissingTitle,
		},
		"effectiveFrom not found": {
			inNoticeRawData: &domain.NoticeRawData{
				PageID:   "test",
				PagePath: "/test",
				Content: `---
title: テストお知らせ
draft: false
effective_until: ` + tomorrow.Format(time.RFC3339) + `
---
これはサンプルです`,
			},

			wantErr: domain.ErrMissingEffectiveFrom,
		},
		"effectiveUntil not found": {
			inNoticeRawData: &domain.NoticeRawData{
				PageID:   "test",
				PagePath: "/test",
				Content: `---
title: テストお知らせ
draft: false
effective_from: ` + yesterday.Format(time.RFC3339) + `
---
これはサンプルです`,
			},

			wantErr: domain.ErrMissingEffectiveUntil,
		},
		"frontMatter not found": {
			inNoticeRawData: &domain.NoticeRawData{
				PageID:   "test",
				PagePath: "/test",
				Content:  "",
			},

			wantErr: domain.ErrMissingFrontMatter,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mockNoticeGetter := noticeGetterFunc(func(ctx context.Context, pagePath string) (*domain.NoticeRawData, error) {
				return tt.inNoticeRawData, nil
			})

			notice, err := domain.FetchNoticeByPath(t.Context(), mockNoticeGetter, "/test")

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
			if tt.wantErr != nil {
				return
			}

			if diff := cmp.Diff(tt.wantNotice, notice,
				cmp.AllowUnexported(domain.Notice{}),        // 非公開フィールドを比較可能にする
				cmpopts.IgnoreFields(domain.Notice{}, "id"), // ID は無視
				cmpopts.EquateApproxTime(time.Second),       // 時刻の誤差を許容
			); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type noticeGetterFunc func(ctx context.Context, pagePath string) (*domain.NoticeRawData, error)

func (f noticeGetterFunc) GetNoticeByPath(ctx context.Context, pagePath string) (*domain.NoticeRawData, error) {
	return f(ctx, pagePath)
}
