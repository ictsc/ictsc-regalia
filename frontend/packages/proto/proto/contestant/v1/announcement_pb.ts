// @generated by protoc-gen-es v2.2.1 with parameter "target=ts"
// @generated from file contestant/v1/announcement.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import { file_buf_validate_validate } from "../../buf/validate/validate_pb";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/announcement.proto.
 */
export const file_contestant_v1_announcement: GenFile = /*@__PURE__*/
  fileDesc("CiBjb250ZXN0YW50L3YxL2Fubm91bmNlbWVudC5wcm90bxINY29udGVzdGFudC52MSLCAQoMQW5ub3VuY2VtZW50EhQKAmlkGAEgASgJQgi6SAVyA5gBGhIhCgpwcm9ibGVtX2lkGAIgASgJQgi6SAVyA5gBGkgAiAEBEhgKBXRpdGxlGAMgASgJQgm6SAZyBBABGGQSGAoEYm9keRgEIAEoCUIKukgHcgUQARjoBxI2CgpjcmVhdGVkX2F0GAUgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcEIGukgDyAEBQg0KC19wcm9ibGVtX2lkIksKF0dldEFubm91bmNlbWVudHNSZXF1ZXN0EiEKCnByb2JsZW1faWQYASABKAlCCLpIBXIDmAEaSACIAQFCDQoLX3Byb2JsZW1faWQiVgoYR2V0QW5ub3VuY2VtZW50c1Jlc3BvbnNlEjoKDWFubm91bmNlbWVudHMYASADKAsyGy5jb250ZXN0YW50LnYxLkFubm91bmNlbWVudEIGukgDyAEBMnoKE0Fubm91bmNlbWVudFNlcnZpY2USYwoQR2V0QW5ub3VuY2VtZW50cxImLmNvbnRlc3RhbnQudjEuR2V0QW5ub3VuY2VtZW50c1JlcXVlc3QaJy5jb250ZXN0YW50LnYxLkdldEFubm91bmNlbWVudHNSZXNwb25zZULBAQoRY29tLmNvbnRlc3RhbnQudjFCEUFubm91bmNlbWVudFByb3RvUAFaRGdpdGh1Yi5jb20vaWN0c2MvaWN0c2Mtb3V0bGFuZHMvYmFja2VuZC9pbnRlcm5hbC9wcm90by9jb250ZXN0YW50L3YxogIDQ1hYqgINQ29udGVzdGFudC5WMcoCDUNvbnRlc3RhbnRcVjHiAhlDb250ZXN0YW50XFYxXEdQQk1ldGFkYXRh6gIOQ29udGVzdGFudDo6VjFiBnByb3RvMw", [file_buf_validate_validate, file_google_protobuf_timestamp]);

/**
 * @generated from message contestant.v1.Announcement
 */
export type Announcement = Message<"contestant.v1.Announcement"> & {
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
 * Describes the message contestant.v1.Announcement.
 * Use `create(AnnouncementSchema)` to create a new message.
 */
export const AnnouncementSchema: GenMessage<Announcement> = /*@__PURE__*/
  messageDesc(file_contestant_v1_announcement, 0);

/**
 * @generated from message contestant.v1.GetAnnouncementsRequest
 */
export type GetAnnouncementsRequest = Message<"contestant.v1.GetAnnouncementsRequest"> & {
  /**
   * @generated from field: optional string problem_id = 1;
   */
  problemId?: string;
};

/**
 * Describes the message contestant.v1.GetAnnouncementsRequest.
 * Use `create(GetAnnouncementsRequestSchema)` to create a new message.
 */
export const GetAnnouncementsRequestSchema: GenMessage<GetAnnouncementsRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_announcement, 1);

/**
 * @generated from message contestant.v1.GetAnnouncementsResponse
 */
export type GetAnnouncementsResponse = Message<"contestant.v1.GetAnnouncementsResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.Announcement announcements = 1;
   */
  announcements: Announcement[];
};

/**
 * Describes the message contestant.v1.GetAnnouncementsResponse.
 * Use `create(GetAnnouncementsResponseSchema)` to create a new message.
 */
export const GetAnnouncementsResponseSchema: GenMessage<GetAnnouncementsResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_announcement, 2);

/**
 * @generated from service contestant.v1.AnnouncementService
 */
export const AnnouncementService: GenService<{
  /**
   * @generated from rpc contestant.v1.AnnouncementService.GetAnnouncements
   */
  getAnnouncements: {
    methodKind: "unary";
    input: typeof GetAnnouncementsRequestSchema;
    output: typeof GetAnnouncementsResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_announcement, 0);

