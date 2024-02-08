// @generated by protoc-gen-connect-es v1.3.0 with parameter "target=ts"
// @generated from file admin/v1/user.proto (package admin.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { DeleteUserRequest, DeleteUserResponse, GetUserRequest, GetUserResponse, GetUsersRequest, GetUsersResponse } from "./user_pb.js";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service admin.v1.UserService
 */
export const UserService = {
  typeName: "admin.v1.UserService",
  methods: {
    /**
     * @generated from rpc admin.v1.UserService.GetUser
     */
    getUser: {
      name: "GetUser",
      I: GetUserRequest,
      O: GetUserResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc admin.v1.UserService.GetUsers
     */
    getUsers: {
      name: "GetUsers",
      I: GetUsersRequest,
      O: GetUsersResponse,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc admin.v1.UserService.DeleteUser
     */
    deleteUser: {
      name: "DeleteUser",
      I: DeleteUserRequest,
      O: DeleteUserResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;

