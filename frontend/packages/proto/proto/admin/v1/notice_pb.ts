// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file admin/v1/notice.proto (package admin.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file admin/v1/notice.proto.
 */
export const file_admin_v1_notice: GenFile = /*@__PURE__*/
  fileDesc("ChVhZG1pbi92MS9ub3RpY2UucHJvdG8SCGFkbWluLnYxImsKBk5vdGljZRIMCgRzbHVnGAEgASgJEg0KBXRpdGxlGAIgASgJEhAKCG1hcmtkb3duGAMgASgJEjIKDmVmZmVjdGl2ZV9mcm9tGAQgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcCIUChJMaXN0Tm90aWNlc1JlcXVlc3QiOAoTTGlzdE5vdGljZXNSZXNwb25zZRIhCgdub3RpY2VzGAEgAygLMhAuYWRtaW4udjEuTm90aWNlIjkKFFVwZGF0ZU5vdGljZXNSZXF1ZXN0EiEKB25vdGljZXMYASADKAsyEC5hZG1pbi52MS5Ob3RpY2UiFwoVVXBkYXRlTm90aWNlc1Jlc3BvbnNlMq0BCg1Ob3RpY2VTZXJ2aWNlEkoKC0xpc3ROb3RpY2VzEhwuYWRtaW4udjEuTGlzdE5vdGljZXNSZXF1ZXN0Gh0uYWRtaW4udjEuTGlzdE5vdGljZXNSZXNwb25zZRJQCg1VcGRhdGVOb3RpY2VzEh4uYWRtaW4udjEuVXBkYXRlTm90aWNlc1JlcXVlc3QaHy5hZG1pbi52MS5VcGRhdGVOb3RpY2VzUmVzcG9uc2VCnwEKDGNvbS5hZG1pbi52MUILTm90aWNlUHJvdG9QAVpBZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1yZWdhbGlhL2JhY2tlbmQvcGtnL3Byb3RvL2FkbWluL3YxO2FkbWludjGiAgNBWFiqAghBZG1pbi5WMcoCCEFkbWluXFYx4gIUQWRtaW5cVjFcR1BCTWV0YWRhdGHqAglBZG1pbjo6VjFiBnByb3RvMw", [file_google_protobuf_timestamp]);

/**
 * @generated from message admin.v1.Notice
 */
export type Notice = Message<"admin.v1.Notice"> & {
  /**
   * @generated from field: string slug = 1;
   */
  slug: string;

  /**
   * @generated from field: string title = 2;
   */
  title: string;

  /**
   * @generated from field: string markdown = 3;
   */
  markdown: string;

  /**
   * @generated from field: google.protobuf.Timestamp effective_from = 4;
   */
  effectiveFrom?: Timestamp;
};

/**
 * Describes the message admin.v1.Notice.
 * Use `create(NoticeSchema)` to create a new message.
 */
export const NoticeSchema: GenMessage<Notice> = /*@__PURE__*/
  messageDesc(file_admin_v1_notice, 0);

/**
 * @generated from message admin.v1.ListNoticesRequest
 */
export type ListNoticesRequest = Message<"admin.v1.ListNoticesRequest"> & {
};

/**
 * Describes the message admin.v1.ListNoticesRequest.
 * Use `create(ListNoticesRequestSchema)` to create a new message.
 */
export const ListNoticesRequestSchema: GenMessage<ListNoticesRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_notice, 1);

/**
 * @generated from message admin.v1.ListNoticesResponse
 */
export type ListNoticesResponse = Message<"admin.v1.ListNoticesResponse"> & {
  /**
   * @generated from field: repeated admin.v1.Notice notices = 1;
   */
  notices: Notice[];
};

/**
 * Describes the message admin.v1.ListNoticesResponse.
 * Use `create(ListNoticesResponseSchema)` to create a new message.
 */
export const ListNoticesResponseSchema: GenMessage<ListNoticesResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_notice, 2);

/**
 * @generated from message admin.v1.UpdateNoticesRequest
 */
export type UpdateNoticesRequest = Message<"admin.v1.UpdateNoticesRequest"> & {
  /**
   * @generated from field: repeated admin.v1.Notice notices = 1;
   */
  notices: Notice[];
};

/**
 * Describes the message admin.v1.UpdateNoticesRequest.
 * Use `create(UpdateNoticesRequestSchema)` to create a new message.
 */
export const UpdateNoticesRequestSchema: GenMessage<UpdateNoticesRequest> = /*@__PURE__*/
  messageDesc(file_admin_v1_notice, 3);

/**
 * @generated from message admin.v1.UpdateNoticesResponse
 */
export type UpdateNoticesResponse = Message<"admin.v1.UpdateNoticesResponse"> & {
};

/**
 * Describes the message admin.v1.UpdateNoticesResponse.
 * Use `create(UpdateNoticesResponseSchema)` to create a new message.
 */
export const UpdateNoticesResponseSchema: GenMessage<UpdateNoticesResponse> = /*@__PURE__*/
  messageDesc(file_admin_v1_notice, 4);

/**
 * @generated from service admin.v1.NoticeService
 */
export const NoticeService: GenService<{
  /**
   * @generated from rpc admin.v1.NoticeService.ListNotices
   */
  listNotices: {
    methodKind: "unary";
    input: typeof ListNoticesRequestSchema;
    output: typeof ListNoticesResponseSchema;
  },
  /**
   * @generated from rpc admin.v1.NoticeService.UpdateNotices
   */
  updateNotices: {
    methodKind: "unary";
    input: typeof UpdateNoticesRequestSchema;
    output: typeof UpdateNoticesResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_admin_v1_notice, 0);

