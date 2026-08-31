# Kubernetes manifests

`base/` contains the shared Deployments and Services. Argo CD deploys
`overlays/dev` from the head of `main` and `overlays/prod` from the latest
stable semantic-version tag.

- A push to `main` publishes the `dev` container images.
- Pushing a stable tag such as `v1.2.3` publishes its semantic-version images
  and moves the `prod` container images to that version.

Argo CD adds the resolved Git commit as a common annotation, so advancing either
revision rolls the Deployments even though the channel image tag is stable.

Each overlay also publishes the contestant and admin frontends through the
cluster's `gateway/external` Gateway. Requests under `/api` are sent to the
backend after removing the prefix. The dev routes use the cluster's
oauth2-proxy ExternalAuth flow; the prod routes do not.

The previous deployment in `ictsc-regalia-release` also contained PostgreSQL,
Valkey, batch workers, and secrets. They are intentionally not copied here until
the rewritten application needs them.
