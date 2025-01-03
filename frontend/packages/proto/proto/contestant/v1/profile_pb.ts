// @generated by protoc-gen-es v2.2.3 with parameter "target=ts"
// @generated from file contestant/v1/profile.proto (package contestant.v1, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file contestant/v1/profile.proto.
 */
export const file_contestant_v1_profile: GenFile = /*@__PURE__*/
  fileDesc("Chtjb250ZXN0YW50L3YxL3Byb2ZpbGUucHJvdG8SDWNvbnRlc3RhbnQudjEiZAoLVGVhbVByb2ZpbGUSDAoEbmFtZRgBIAEoCRIUCgxvcmdhbml6YXRpb24YAiABKAkSMQoHbWVtYmVycxgDIAMoCzIgLmNvbnRlc3RhbnQudjEuQ29udGVzdGFudFByb2ZpbGUiUgoRQ29udGVzdGFudFByb2ZpbGUSDAoEbmFtZRgBIAEoCRIUCgxkaXNwbGF5X25hbWUYAiABKAkSGQoRc2VsZl9pbnRyb2R1Y3Rpb24YAyABKAkiEgoQTGlzdFRlYW1zUmVxdWVzdCI+ChFMaXN0VGVhbXNSZXNwb25zZRIpCgV0ZWFtcxgBIAMoCzIaLmNvbnRlc3RhbnQudjEuVGVhbVByb2ZpbGUiSQoUVXBkYXRlUHJvZmlsZVJlcXVlc3QSMQoHcHJvZmlsZRgBIAEoCzIgLmNvbnRlc3RhbnQudjEuQ29udGVzdGFudFByb2ZpbGUiSgoVVXBkYXRlUHJvZmlsZVJlc3BvbnNlEjEKB3Byb2ZpbGUYASABKAsyIC5jb250ZXN0YW50LnYxLkNvbnRlc3RhbnRQcm9maWxlMrwBCg5Qcm9maWxlU2VydmljZRJOCglMaXN0VGVhbXMSHy5jb250ZXN0YW50LnYxLkxpc3RUZWFtc1JlcXVlc3QaIC5jb250ZXN0YW50LnYxLkxpc3RUZWFtc1Jlc3BvbnNlEloKDVVwZGF0ZVByb2ZpbGUSIy5jb250ZXN0YW50LnYxLlVwZGF0ZVByb2ZpbGVSZXF1ZXN0GiQuY29udGVzdGFudC52MS5VcGRhdGVQcm9maWxlUmVzcG9uc2VCwwEKEWNvbS5jb250ZXN0YW50LnYxQgxQcm9maWxlUHJvdG9QAVpLZ2l0aHViLmNvbS9pY3RzYy9pY3RzYy1yZWdhbGlhL2JhY2tlbmQvcGtnL3Byb3RvL2NvbnRlc3RhbnQvdjE7Y29udGVzdGFudHYxogIDQ1hYqgINQ29udGVzdGFudC5WMcoCDUNvbnRlc3RhbnRcVjHiAhlDb250ZXN0YW50XFYxXEdQQk1ldGFkYXRh6gIOQ29udGVzdGFudDo6VjFiBnByb3RvMw");

/**
 * @generated from message contestant.v1.TeamProfile
 */
export type TeamProfile = Message<"contestant.v1.TeamProfile"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: string organization = 2;
   */
  organization: string;

  /**
   * @generated from field: repeated contestant.v1.ContestantProfile members = 3;
   */
  members: ContestantProfile[];
};

/**
 * Describes the message contestant.v1.TeamProfile.
 * Use `create(TeamProfileSchema)` to create a new message.
 */
export const TeamProfileSchema: GenMessage<TeamProfile> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 0);

/**
 * @generated from message contestant.v1.ContestantProfile
 */
export type ContestantProfile = Message<"contestant.v1.ContestantProfile"> & {
  /**
   * @generated from field: string name = 1;
   */
  name: string;

  /**
   * @generated from field: string display_name = 2;
   */
  displayName: string;

  /**
   * @generated from field: string self_introduction = 3;
   */
  selfIntroduction: string;
};

/**
 * Describes the message contestant.v1.ContestantProfile.
 * Use `create(ContestantProfileSchema)` to create a new message.
 */
export const ContestantProfileSchema: GenMessage<ContestantProfile> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 1);

/**
 * @generated from message contestant.v1.ListTeamsRequest
 */
export type ListTeamsRequest = Message<"contestant.v1.ListTeamsRequest"> & {
};

/**
 * Describes the message contestant.v1.ListTeamsRequest.
 * Use `create(ListTeamsRequestSchema)` to create a new message.
 */
export const ListTeamsRequestSchema: GenMessage<ListTeamsRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 2);

/**
 * @generated from message contestant.v1.ListTeamsResponse
 */
export type ListTeamsResponse = Message<"contestant.v1.ListTeamsResponse"> & {
  /**
   * @generated from field: repeated contestant.v1.TeamProfile teams = 1;
   */
  teams: TeamProfile[];
};

/**
 * Describes the message contestant.v1.ListTeamsResponse.
 * Use `create(ListTeamsResponseSchema)` to create a new message.
 */
export const ListTeamsResponseSchema: GenMessage<ListTeamsResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 3);

/**
 * @generated from message contestant.v1.UpdateProfileRequest
 */
export type UpdateProfileRequest = Message<"contestant.v1.UpdateProfileRequest"> & {
  /**
   * @generated from field: contestant.v1.ContestantProfile profile = 1;
   */
  profile?: ContestantProfile;
};

/**
 * Describes the message contestant.v1.UpdateProfileRequest.
 * Use `create(UpdateProfileRequestSchema)` to create a new message.
 */
export const UpdateProfileRequestSchema: GenMessage<UpdateProfileRequest> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 4);

/**
 * @generated from message contestant.v1.UpdateProfileResponse
 */
export type UpdateProfileResponse = Message<"contestant.v1.UpdateProfileResponse"> & {
  /**
   * @generated from field: contestant.v1.ContestantProfile profile = 1;
   */
  profile?: ContestantProfile;
};

/**
 * Describes the message contestant.v1.UpdateProfileResponse.
 * Use `create(UpdateProfileResponseSchema)` to create a new message.
 */
export const UpdateProfileResponseSchema: GenMessage<UpdateProfileResponse> = /*@__PURE__*/
  messageDesc(file_contestant_v1_profile, 5);

/**
 * @generated from service contestant.v1.ProfileService
 */
export const ProfileService: GenService<{
  /**
   * @generated from rpc contestant.v1.ProfileService.ListTeams
   */
  listTeams: {
    methodKind: "unary";
    input: typeof ListTeamsRequestSchema;
    output: typeof ListTeamsResponseSchema;
  },
  /**
   * @generated from rpc contestant.v1.ProfileService.UpdateProfile
   */
  updateProfile: {
    methodKind: "unary";
    input: typeof UpdateProfileRequestSchema;
    output: typeof UpdateProfileResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_contestant_v1_profile, 0);
