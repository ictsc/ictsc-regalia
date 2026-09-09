import { HttpResponse, http } from "msw";

export { HttpResponse, http };

export function problem(
  status: number,
  code: string,
  detail: string,
  headers?: Record<string, string>,
) {
  return HttpResponse.json(
    {
      type: "about:blank",
      title: status === 429 ? "Too Many Requests" : "Request failed",
      status,
      detail,
      code,
      errors: [],
    },
    { status, headers },
  );
}
