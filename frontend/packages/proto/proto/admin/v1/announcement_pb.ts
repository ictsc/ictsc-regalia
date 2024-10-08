// @generated by protoc-gen-es v2.2.0 with parameter "target=ts"
// @generated from file admin/v1/announcement.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/announcement.proto.
 */
export const file_admin_v1_announcement: GenFile = /*@__PURE__*/
  fileDesc("ChthZG1pbi92MS9hbm5vdW5jZW1lbnQucHJvdG8SCGFkbWluLnYxIsIBCgxBbm5vdW5jZW1lbnQSFAoCaWQYASABKAlCCLpIBXIDmAEaEiEKCnByb2JsZW1faWQYAiABKAlCCLpIBXIDmAEaSACIAQESGAoFdGl0bGUYAyABKAlCCbpIBnIEEAEYZBIYCgRib2R5GAQgASgJQgq6SAdyBRABGOgHEjYKCmNyZWF0ZWRfYXQYBSABKAsyGi5nb29nbGUucHJvdG9idWYuVGltZXN0YW1wQga6SAPIAQFCDQoLX3Byb2JsZW1faWQiSwoXR2V0QW5ub3VuY2VtZW50c1JlcXVlc3QSIQoKcHJvYmxlbV9pZBgBIAEoCUIIukgFcgOYARpIAIgBAUINCgtfcHJvYmxlbV9pZCJRChhHZXRBbm5vdW5jZW1lbnRzUmVzcG9uc2USNQoNYW5ub3VuY2VtZW50cxgBIAMoCzIWLmFkbWluLnYxLkFubm91bmNlbWVudEIGukgDyAEBIrMBChhQYXRjaEFubm91bmNlbWVudFJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaEiEKCnByb2JsZW1faWQYAiABKAlCCLpIBXIDmAEaSACIAQESHQoFdGl0bGUYAyABKAlCCbpIBnIEEAEYZEgBiAEBEh0KBGJvZHkYBCABKAlCCrpIB3IFEAEY6AdIAogBAUINCgtfcHJvYmxlbV9pZEIICgZfdGl0bGVCBwoFX2JvZHkiUQoZUGF0Y2hBbm5vdW5jZW1lbnRSZXNwb25zZRI0Cgxhbm5vdW5jZW1lbnQYASABKAsyFi5hZG1pbi52MS5Bbm5vdW5jZW1lbnRCBrpIA8gBASJrChdQb3N0QW5ub3VuY2VtZW50UmVxdWVzdBIcCgpwcm9ibGVtX2lkGAEgASgJQgi6SAVyA5gBGhIYCgV0aXRsZRgCIAEoCUIJukgGcgQQARhkEhgKBGJvZHkYAyABKAlCCrpIB3IFEAEY6AciUAoYUG9zdEFubm91bmNlbWVudFJlc3BvbnNlEjQKDGFubm91bmNlbWVudBgBIAEoCzIWLmFkbWluLnYxLkFubm91bmNlbWVudEIGukgDyAEBIjEKGURlbGV0ZUFubm91bmNlbWVudFJlcXVlc3QSFAoCaWQYASABKAlCCLpIBXIDmAEaIhwKGkRlbGV0ZUFubm91bmNlbWVudFJlc3BvbnNlMooDChNBbm5vdW5jZW1lbnRTZXJ2aWNlElkKEEdldEFubm91bmNlbWVudHMSIS5hZG1pbi52MS5HZXRBbm5vdW5jZW1lbnRzUmVxdWVzdBoiLmFkbWluLnYxLkdldEFubm91bmNlbWVudHNSZXNwb25zZRJcChFQYXRjaEFubm91bmNlbWVudBIiLmFkbWluLnYxLlBhdGNoQW5ub3VuY2VtZW50UmVxdWVzdBojLmFkbWluLnYxLlBhdGNoQW5ub3VuY2VtZW50UmVzcG9uc2USWQoQUG9zdEFubm91bmNlbWVudBIhLmFkbWluLnYxLlBvc3RBbm5vdW5jZW1lbnRSZXF1ZXN0GiIuYWRtaW4udjEuUG9zdEFubm91bmNlbWVudFJlc3BvbnNlEl8KEkRlbGV0ZUFubm91bmNlbWVudBIjLmFkbWluLnYxLkRlbGV0ZUFubm91bmNlbWVudFJlcXVlc3QaJC5hZG1pbi52MS5EZWxldGVBbm5vdW5jZW1lbnRSZXNwb25zZUKjAQoMY29tLmFkbWluLnYxQhFBbm5vdW5jZW1lbnRQcm90b1ABWj9naXRodWIuY29tL2ljdHNjL2ljdHNjLW91dGxhbmRzL2JhY2tlbmQvaW50ZXJuYWwvcHJvdG8vYWRtaW4vdjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_buf_validate_validate, file_google_protobuf_timestamp]);

/**
 * @generated from message admin.v1.Announcement
 */
export type Announcement = Message<"admin.v1.Announcement"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: optional string problem_id = 2;
   */
  problemId?: string;

  /**
   * @generated from field: string title = 3;
   */
  title: string;

  /**
   * @generated from field: string body = 4;
   */
  body: string;

  /**
   * @generated from field: google.protobuf.Timestamp created_at = 5;
   */
  createdAt?: Timestamp;
};

/**
 * Describes the message admin.v1.Announcement.
 * Use `create(AnnouncementSchema)` to create a new message.
 */
export const AnnouncementSchema: GenMessage<Announcement> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 0);

/**
 * @generated from message admin.v1.GetAnnouncementsRequest
 */
export type GetAnnouncementsRequest = Message<"admin.v1.GetAnnouncementsRequest"> & {
  /**
   * @generated from field: optional string problem_id = 1;
   */
  problemId?: string;
};

/**
 * Describes the message admin.v1.GetAnnouncementsRequest.
 * Use `create(GetAnnouncementsRequestSchema)` to create a new message.
 */
export const GetAnnouncementsRequestSchema: GenMessage<GetAnnouncementsRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 1);

/**
 * @generated from message admin.v1.GetAnnouncementsResponse
 */
export type GetAnnouncementsResponse = Message<"admin.v1.GetAnnouncementsResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Announcement announcements = 1;
   */
  announcements: Announcement[];
};

/**
 * Describes the message admin.v1.GetAnnouncementsResponse.
 * Use `create(GetAnnouncementsResponseSchema)` to create a new message.
 */
export const GetAnnouncementsResponseSchema: GenMessage<GetAnnouncementsResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 2);

/**
 * @generated from message admin.v1.PatchAnnouncementRequest
 */
export type PatchAnnouncementRequest = Message<"admin.v1.PatchAnnouncementRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: optional string problem_id = 2;
   */
  problemId?: string;

  /**
   * @generated from field: optional string title = 3;
   */
  title?: string;

  /**
   * @generated from field: optional string body = 4;
   */
  body?: string;
};

/**
 * Describes the message admin.v1.PatchAnnouncementRequest.
 * Use `create(PatchAnnouncementRequestSchema)` to create a new message.
 */
export const PatchAnnouncementRequestSchema: GenMessage<PatchAnnouncementRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 3);

/**
 * @generated from message admin.v1.PatchAnnouncementResponse
 */
export type PatchAnnouncementResponse = Message<"admin.v1.PatchAnnouncementResponse"> & {
  /**
   * @generated from field: admin.v1.Announcement announcement = 1;
   */
  announcement?: Announcement;
};

/**
 * Describes the message admin.v1.PatchAnnouncementResponse.
 * Use `create(PatchAnnouncementResponseSchema)` to create a new message.
 */
export const PatchAnnouncementResponseSchema: GenMessage<PatchAnnouncementResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 4);

/**
 * @generated from message admin.v1.PostAnnouncementRequest
 */
export type PostAnnouncementRequest = Message<"admin.v1.PostAnnouncementRequest"> & {
  /**
   * @generated from field: string problem_id = 1;
   */
  problemId: string;

  /**
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * @generated from field: string body = 3;
   */
  body: string;
};

/**
 * Describes the message admin.v1.PostAnnouncementRequest.
 * Use `create(PostAnnouncementRequestSchema)` to create a new message.
 */
export const PostAnnouncementRequestSchema: GenMessage<PostAnnouncementRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 5);

/**
 * @generated from message admin.v1.PostAnnouncementResponse
 */
export type PostAnnouncementResponse = Message<"admin.v1.PostAnnouncementResponse"> & {
  /**
   * @generated from field: admin.v1.Announcement announcement = 1;
   */
  announcement?: Announcement;
};

/**
 * Describes the message admin.v1.PostAnnouncementResponse.
 * Use `create(PostAnnouncementResponseSchema)` to create a new message.
 */
export const PostAnnouncementResponseSchema: GenMessage<PostAnnouncementResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 6);

/**
 * @generated from message admin.v1.DeleteAnnouncementRequest
 */
export type DeleteAnnouncementRequest = Message<"admin.v1.DeleteAnnouncementRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message admin.v1.DeleteAnnouncementRequest.
 * Use `create(DeleteAnnouncementRequestSchema)` to create a new message.
 */
export const DeleteAnnouncementRequestSchema: GenMessage<DeleteAnnouncementRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 7);

/**
 * @generated from message admin.v1.DeleteAnnouncementResponse
 */
export type DeleteAnnouncementResponse = Message<"admin.v1.DeleteAnnouncementResponse"> & {
};

/**
 * Describes the message admin.v1.DeleteAnnouncementResponse.
 * Use `create(DeleteAnnouncementResponseSchema)` to create a new message.
 */
export const DeleteAnnouncementResponseSchema: GenMessage<DeleteAnnouncementResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_announcement, 8);

/**
 * @generated from service admin.v1.AnnouncementService
 */
export const AnnouncementService: GenService<{
  /**
   * @generated from rpc admin.v1.AnnouncementService.GetAnnouncements
   */
  getAnnouncements: {
    methodKind: "unary";
    input: typeof GetAnnouncementsRequestSchema;
    output: typeof GetAnnouncementsResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.AnnouncementService.PatchAnnouncement
   */
  patchAnnouncement: {
    methodKind: "unary";
    input: typeof PatchAnnouncementRequestSchema;
    output: typeof PatchAnnouncementResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.AnnouncementService.PostAnnouncement
   */
  postAnnouncement: {
    methodKind: "unary";
    input: typeof PostAnnouncementRequestSchema;
    output: typeof PostAnnouncementResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.AnnouncementService.DeleteAnnouncement
   */
  deleteAnnouncement: {
    methodKind: "unary";
    input: typeof DeleteAnnouncementRequestSchema;
    output: typeof DeleteAnnouncementResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_announcement, 0);

