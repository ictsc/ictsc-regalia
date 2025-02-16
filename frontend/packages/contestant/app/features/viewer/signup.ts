export type SignUpRequest = {
  invitationCode: string;
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
  let error = false;
  if (request.invitationCode === "") {
    result.invitationCodeError = "required";
    error = true;
  }
  if (request.name === "") {
    result.nameError = "required";
    error = true;
  }
  if (request.displayName === "") {
    result.displayNameError = "required";
    error = true;
  }
  if (error) {
    return result;
  }

  const url = new URL(baseURL ?? window.location.href);
  url.pathname = "/api/auth/signup";
  const resp = await fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      invitation_code: request.invitationCode,
      name: request.name,
      display_name: request.displayName,
    }),
  });
  if (resp.status === 429) {
    return { error: "rate_limit" };
  }
  if (resp.status === 400) {
    result.error = "invalid";
  }
  const body: unknown = await resp.json();
  if (typeof body !== "object" || body === null) {
    throw new Error("unexpected response");
  }
  let codes: unknown;
  if ("codes" in body) {
    codes = body.codes;
  } else {
    codes = [];
  }

  if (typeof codes !== "object" || codes === null || !Array.isArray(codes)) {
    throw new Error("unexpected response");
  }
  for (const code of codes) {
    switch (code) {
      case "invalid_invitation_code":
        result.invitationCodeError = "invalid";
        break;
      case "team_is_full":
        result.invitationCodeError = "team_full";
        break;
      case "invalid_name":
        result.nameError = "invalid";
        break;
      case "duplicate_name":
        result.nameError = "duplicate";
        break;
      case "invalid_display_name":
        result.displayNameError = "invalid";
        break;
      default:
        break;
    }
  }
  return result;
}
