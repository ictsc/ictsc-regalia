export { api, createApiClient, type ApiClient } from "./client";
export {
  ApiError,
  expectData,
  expectNoContent,
  isApiError,
  type ProblemDetails,
} from "./error";
export { toCamelCase, toSnakeCase } from "./mapper";
export {
  subscribeAdminAnswerEvents,
  subscribeAdminDeploymentEvents,
  subscribeContestantAnswerEvents,
  subscribeContestantDeploymentEvents,
  type AdminAnswerChange,
  type AnswerStreamMessage,
  type ContestantAnswerChange,
  type ContestantDeployment,
  type AdminDeployment,
  type ContestantDeploymentStreamMessage,
  type AdminDeploymentStreamMessage,
} from "./sse";
export {
  createRefetchScheduler,
  type RefetchScheduler,
  type RefetchSchedulerOptions,
} from "./refetch";
export type { paths, components, operations } from "./schema";

export { impersonateContestant } from "./impersonation";
