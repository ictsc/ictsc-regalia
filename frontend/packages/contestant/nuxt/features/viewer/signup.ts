import { api, ApiError, expectNoContent } from "@ictsc/api";

export type SignUpRequest = {
  invitationCode?: string;
  name: string;
  displayName: string;
};

export type SignUpResponse = {
  error?: "rate_limit" | "invalid" | "unknown";
  invitationCodeError?: "required" | "invalid" | "team_full";
  nameError?: "required" | "invalid" | "duplicate";
  displayNameError?: "required" | "invalid";
};

export async function signUp(
  request: SignUpRequest,
  baseURL?: string,
): Promise<SignUpResponse> {
  const result: SignUpResponse = {};
  if (request.invitationCode === "") {
    result.invitationCodeError = "required";
    result.error = "invalid";
  }
  if (request.name === "") {
    result.nameError = "required";
    result.error = "invalid";
  }
  if (request.displayName === "") {
    result.displayNameError = "required";
    result.error = "invalid";
  }
  if (result.error != null) {
    return result;
  }

  try {
    const client =
      baseURL == null
        ? api
        : (await import("@ictsc/api")).createApiClient(baseURL);
    expectNoContent(
      await client.POST("/api/v1/auth/signup", {
        body: {
          invitation_code: request.invitationCode,
          name: request.name,
          display_name: request.displayName,
        },
      }),
    );
    return result;
  } catch (error) {
    if (!(error instanceof ApiError)) throw error;
    if (error.status === 429) return { error: "rate_limit" };
    result.error =
      error.status === 409 || error.status === 422 ? "invalid" : "unknown";
    const code = error.code;
    if (code === "invalid_invitation_code" || code === "invitation_expired") {
      result.invitationCodeError = "invalid";
    } else if (code === "team_full") {
      result.invitationCodeError = "team_full";
    } else if (code === "contestant_name_conflict") {
      result.nameError = "duplicate";
    }
    for (const detail of error.problem?.errors ?? []) {
      const location = detail.location;
      if (location.includes("invitation_code"))
        result.invitationCodeError ??= "invalid";
      if (location.endsWith("name")) result.nameError ??= "invalid";
      if (location.includes("display_name"))
        result.displayNameError ??= "invalid";
    }
    return result;
  }
}
