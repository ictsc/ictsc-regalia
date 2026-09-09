import type { components } from "./schema";

export type ProblemDetails = components["schemas"]["ProblemDetails"];

type ApiResult<T> = {
  data?: T;
  error?: unknown;
  response: Response;
};

export class ApiError extends Error {
  readonly status: number;
  readonly code?: string;
  readonly problem?: ProblemDetails;
  readonly retryAfterSeconds?: number;

  constructor(
    message: string,
    options: {
      status: number;
      code?: string;
      problem?: ProblemDetails;
      retryAfterSeconds?: number;
    },
  ) {
    super(message);
    this.name = "ApiError";
    this.status = options.status;
    this.code = options.code;
    this.problem = options.problem;
    this.retryAfterSeconds = options.retryAfterSeconds;
  }

  static fromResponse(response: Response, body?: unknown): ApiError {
    const problem = isProblemDetails(body) ? body : undefined;
    const retryAfterSeconds = parseRetryAfter(
      response.headers.get("Retry-After"),
    );
    return new ApiError(
      problem?.detail ?? problem?.title ?? `HTTP ${response.status}`,
      {
        status: response.status,
        code: problem?.code,
        problem,
        retryAfterSeconds,
      },
    );
  }
}

export function expectData<T>(result: ApiResult<T>): T {
  if (!result.response.ok || result.error !== undefined) {
    throw ApiError.fromResponse(result.response, result.error);
  }
  if (result.data === undefined) {
    throw new ApiError("API response body is missing", {
      status: result.response.status,
      code: "response_body_missing",
    });
  }
  return result.data;
}

export function expectNoContent(result: ApiResult<unknown>): void {
  if (!result.response.ok || result.error !== undefined) {
    throw ApiError.fromResponse(result.response, result.error);
  }
}

export function isApiError(error: unknown, status?: number): error is ApiError {
  return (
    error instanceof ApiError &&
    (status === undefined || error.status === status)
  );
}

function isProblemDetails(value: unknown): value is ProblemDetails {
  if (value == null || typeof value !== "object") return false;
  const candidate = value as Record<string, unknown>;
  return (
    typeof candidate.title === "string" && typeof candidate.status === "number"
  );
}

function parseRetryAfter(value: string | null): number | undefined {
  if (value == null) return undefined;
  const seconds = Number(value);
  if (Number.isFinite(seconds) && seconds >= 0) return Math.ceil(seconds);
  const date = Date.parse(value);
  if (Number.isNaN(date)) return undefined;
  return Math.max(0, Math.ceil((date - Date.now()) / 1000));
}
