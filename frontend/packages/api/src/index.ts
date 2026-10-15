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
  subscribeAdminDeploymentEvents,
  subscribeContestantDeploymentEvents,
  type ContestantDeployment,
  type AdminDeployment,
  type ContestantDeploymentStreamMessage,
  type AdminDeploymentStreamMessage,
} from "./sse";
export type { paths, components, operations } from "./schema";

export { impersonateContestant } from "./impersonation";
