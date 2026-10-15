package postgres

import (
	"context"
	"time"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
)

func (s *Store) UpsertWebPushSubscription(ctx context.Context, subscription core.WebPushSubscription) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO web_push_subscriptions(contestant_name,endpoint,p256dh,auth,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$5)
		ON CONFLICT(endpoint) DO UPDATE SET
			contestant_name=EXCLUDED.contestant_name,
			p256dh=EXCLUDED.p256dh,
			auth=EXCLUDED.auth,
			updated_at=EXCLUDED.updated_at`,
		subscription.ContestantName, subscription.Endpoint, subscription.P256DH, subscription.Auth, subscription.UpdatedAt)
	return dbError(err)
}

func (s *Store) DeleteWebPushSubscription(ctx context.Context, contestantName, endpoint string) error {
	_, err := s.pool.Exec(ctx, `
		DELETE FROM web_push_subscriptions
		WHERE contestant_name=$1 AND endpoint=$2`, contestantName, endpoint)
	return dbError(err)
}

func (s *Store) ListWebPushSubscriptions(ctx context.Context) ([]core.WebPushSubscription, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT contestant_name,endpoint,p256dh,auth,created_at,updated_at
		FROM web_push_subscriptions ORDER BY created_at,endpoint`)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()
	result := make([]core.WebPushSubscription, 0)
	for rows.Next() {
		var subscription core.WebPushSubscription
		if err := rows.Scan(&subscription.ContestantName, &subscription.Endpoint, &subscription.P256DH,
			&subscription.Auth, &subscription.CreatedAt, &subscription.UpdatedAt); err != nil {
			return nil, dbError(err)
		}
		result = append(result, subscription)
	}
	return result, dbError(rows.Err())
}

func (s *Store) ClaimAnnouncementPush(ctx context.Context, announcementSlug, endpoint string, at time.Time) (bool, error) {
	result, err := s.pool.Exec(ctx, `
		INSERT INTO announcement_push_deliveries(announcement_slug,endpoint,claimed_at)
		VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, announcementSlug, endpoint, at)
	if err != nil {
		return false, dbError(err)
	}
	return result.RowsAffected() == 1, nil
}

func (s *Store) ReleaseAnnouncementPush(ctx context.Context, announcementSlug, endpoint string) error {
	_, err := s.pool.Exec(ctx, `
		DELETE FROM announcement_push_deliveries
		WHERE announcement_slug=$1 AND endpoint=$2`, announcementSlug, endpoint)
	return dbError(err)
}
