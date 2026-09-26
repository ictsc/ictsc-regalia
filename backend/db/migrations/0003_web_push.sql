BEGIN;

CREATE TABLE web_push_subscriptions (
    endpoint TEXT PRIMARY KEY CHECK (btrim(endpoint) <> ''),
    contestant_name VARCHAR(32) NOT NULL REFERENCES contestants(name) ON DELETE CASCADE,
    p256dh TEXT NOT NULL CHECK (btrim(p256dh) <> ''),
    auth TEXT NOT NULL CHECK (btrim(auth) <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX web_push_subscriptions_contestant_idx
    ON web_push_subscriptions(contestant_name);

CREATE TABLE announcement_push_deliveries (
    announcement_slug VARCHAR(255) NOT NULL,
    endpoint TEXT NOT NULL REFERENCES web_push_subscriptions(endpoint) ON DELETE CASCADE,
    claimed_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (announcement_slug, endpoint)
);

COMMIT;
