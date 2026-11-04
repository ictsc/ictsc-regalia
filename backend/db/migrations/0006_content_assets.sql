BEGIN;

-- Images from the content repository's images/ directory. Image bytes are
-- stored once per Git blob SHA, and each snapshot maps its file names to
-- blobs, so refreshing content does not copy unchanged images.
CREATE TABLE content_blobs (
    sha VARCHAR(64) PRIMARY KEY CHECK (sha ~ '^[0-9a-f]{40,64}$'),
    data BYTEA NOT NULL
);

CREATE TABLE content_assets (
    commit_sha VARCHAR(64) NOT NULL REFERENCES content_snapshots(commit_sha) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (name ~ '^[A-Za-z0-9._-]+$'),
    content_type TEXT NOT NULL,
    blob_sha VARCHAR(64) NOT NULL REFERENCES content_blobs(sha) ON DELETE RESTRICT,
    PRIMARY KEY (commit_sha, name)
);

COMMIT;
