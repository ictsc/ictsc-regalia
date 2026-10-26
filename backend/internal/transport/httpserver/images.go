package httpserver

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/ictsc/ictsc-regalia/backend/internal/core"
	api "github.com/ictsc/ictsc-regalia/backend/internal/transport/api"
)

// imageCacheControl lets browsers reuse images while staying private to the
// signed-in viewer; the ETag lets them revalidate cheaply afterwards.
var (
	imageCacheControl = "private, max-age=300"
	imageNoSniff      = "nosniff"
)

func (h *Handler) GetContestantImage(ctx context.Context, request api.GetContestantImageRequestObject) (api.GetContestantImageResponseObject, error) {
	if _, err := h.currentContestant(ctx); err != nil {
		return nil, err
	}
	snapshot, err := h.service.Store.ActiveContent(ctx)
	if err != nil {
		return nil, err
	}
	// Images of problems that have not started stay hidden like the problems.
	if _, visible := core.ContestantVisibleImages(snapshot, h.service.Now())[request.ImageName]; !visible {
		return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Image was not found")
	}
	asset, err := h.service.Store.GetContentAsset(ctx, snapshot.CommitSHA, request.ImageName)
	if err != nil {
		return nil, err
	}
	etag := imageETag(asset)
	if etagMatches(request.Params.IfNoneMatch, etag) {
		return api.GetContestantImage304Response{Headers: api.GetContestantImage304ResponseHeaders{ETag: &etag}}, nil
	}
	body, length := bytes.NewReader(asset.Data), int64(len(asset.Data))
	headers := api.GetContestantImage200ResponseHeaders{CacheControl: &imageCacheControl, ETag: &etag, XContentTypeOptions: &imageNoSniff}
	switch asset.ContentType {
	case "image/png":
		return api.GetContestantImage200ImagepngResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/jpeg":
		return api.GetContestantImage200ImagejpegResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/gif":
		return api.GetContestantImage200ImagegifResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/webp":
		return api.GetContestantImage200ImagewebpResponse{Body: body, ContentLength: length, Headers: headers}, nil
	}
	return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Image was not found")
}

func (h *Handler) GetAdminImage(ctx context.Context, request api.GetAdminImageRequestObject) (api.GetAdminImageResponseObject, error) {
	commit := ""
	if request.Params.Commit != nil {
		commit = string(*request.Params.Commit)
	} else {
		snapshot, err := h.service.Store.ActiveContent(ctx)
		if err != nil {
			return nil, err
		}
		commit = snapshot.CommitSHA
	}
	asset, err := h.service.Store.GetContentAsset(ctx, commit, request.ImageName)
	if err != nil {
		return nil, err
	}
	etag := imageETag(asset)
	if etagMatches(request.Params.IfNoneMatch, etag) {
		return api.GetAdminImage304Response{Headers: api.GetAdminImage304ResponseHeaders{ETag: &etag}}, nil
	}
	body, length := bytes.NewReader(asset.Data), int64(len(asset.Data))
	headers := api.GetAdminImage200ResponseHeaders{CacheControl: &imageCacheControl, ETag: &etag, XContentTypeOptions: &imageNoSniff}
	switch asset.ContentType {
	case "image/png":
		return api.GetAdminImage200ImagepngResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/jpeg":
		return api.GetAdminImage200ImagejpegResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/gif":
		return api.GetAdminImage200ImagegifResponse{Body: body, ContentLength: length, Headers: headers}, nil
	case "image/webp":
		return api.GetAdminImage200ImagewebpResponse{Body: body, ContentLength: length, Headers: headers}, nil
	}
	return nil, core.NewError(http.StatusNotFound, "resource_not_found", "Image was not found")
}

// imageETag uses the Git blob SHA, which changes exactly when the bytes do.
func imageETag(asset core.ContentAsset) string {
	return `"` + asset.SHA + `"`
}

func etagMatches(header *api.IfNoneMatchHeader, etag string) bool {
	if header == nil || etag == `""` {
		return false
	}
	for _, candidate := range strings.Split(string(*header), ",") {
		candidate = strings.TrimPrefix(strings.TrimSpace(candidate), "W/")
		if candidate == etag || candidate == "*" {
			return true
		}
	}
	return false
}
