# OpenAPI acceptance matrix

`backend/openapi.json` の61 operationを正本として、transport・永続化/adapter・frontend・検証の対応を一覧化する。strict server interfaceのコンパイルと `server_contract_test.go` の全operation登録検査を全行の共通ゲートとし、UIを持たないcallback/healthも未利用ではなく明示的なmachine/infra consumerとして扱う。

| operationId | Method / path | Go transport | Domain / persistence | Consumer | Verification |
|---|---|---|---|---|---|
| `getHealth` | `GET /api/v1/health` | strict handler ✓ | service + store/adapter ✓ | infra healthcheck | contract + smoke |
| `getViewer` | `GET /api/v1/viewer` | strict handler ✓ | service + store/adapter ✓ | Contestant root | contract + auth/security unit |
| `startDiscordAuthentication` | `GET /api/v1/auth/discord` | strict handler ✓ | service + store/adapter ✓ | Contestant auth redirect | contract + auth/security unit |
| `completeDiscordAuthentication` | `GET /api/v1/auth/discord/callback` | strict handler ✓ | service + store/adapter ✓ | Contestant auth redirect | contract + auth/security unit |
| `signUpContestant` | `POST /api/v1/auth/signup` | strict handler ✓ | service + store/adapter ✓ | Contestant /signup | contract + auth/security unit |
| `signOutContestant` | `POST /api/v1/auth/signout` | strict handler ✓ | service + store/adapter ✓ | Contestant account menu | contract + auth/security unit |
| `getAdminViewer` | `GET /api/v1/admin/viewer` | strict handler ✓ | service + store/adapter ✓ | Admin root | contract + auth/security unit |
| `startAdminDiscordAuthentication` | `GET /api/v1/admin/auth/discord` | strict handler ✓ | service + store/adapter ✓ | Admin auth redirect | contract + auth/security unit |
| `completeAdminDiscordAuthentication` | `GET /api/v1/admin/auth/discord/callback` | strict handler ✓ | service + store/adapter ✓ | Admin auth redirect | contract + auth/security unit |
| `signOutAdmin` | `POST /api/v1/admin/auth/signout` | strict handler ✓ | service + store/adapter ✓ | Admin root | contract + auth/security unit |
| `createContestantImpersonation` | `POST /api/v1/admin/impersonations` | strict handler ✓ | service + store/adapter ✓ | Admin /contestants | contract + auth/security unit |
| `getContestantProfile` | `GET /api/v1/contestant/profile` | strict handler ✓ | service + store/adapter ✓ | Contestant /profile | contract + smoke |
| `updateContestantProfile` | `PATCH /api/v1/contestant/profile` | strict handler ✓ | service + store/adapter ✓ | Contestant /profile | contract + smoke |
| `listContestantTeams` | `GET /api/v1/contestant/teams` | strict handler ✓ | service + store/adapter ✓ | Contestant /teams | contract + service/PostgreSQL + E2E |
| `listContestantSections` | `GET /api/v1/contestant/sections` | strict handler ✓ | service + store/adapter ✓ | Contestant dashboard | contract + content unit/integration + E2E |
| `getContestantSchedule` | `GET /api/v1/contestant/schedule` | strict handler ✓ | service + store/adapter ✓ | Contestant dashboard | contract + service/PostgreSQL + E2E |
| `listContestantProblems` | `GET /api/v1/contestant/problems` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `getContestantProblem` | `GET /api/v1/contestant/problems/{problem_code}` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `listContestantAnswers` | `GET /api/v1/contestant/problems/{problem_code}/answers` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `submitContestantAnswer` | `POST /api/v1/contestant/problems/{problem_code}/answers` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `getContestantAnswer` | `GET /api/v1/contestant/problems/{problem_code}/answers/{answer_number}` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `listContestantDeployments` | `GET /api/v1/contestant/problems/{problem_code}/deployments` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `createContestantDeployment` | `POST /api/v1/contestant/problems/{problem_code}/deployments` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `streamContestantDeployments` | `GET /api/v1/contestant/problems/{problem_code}/deployments/stream` | strict handler ✓ | service + store/adapter ✓ | Contestant /problems | contract + content unit/integration + E2E |
| `listContestantAnnouncements` | `GET /api/v1/contestant/announcements` | strict handler ✓ | service + store/adapter ✓ | Contestant /announces | contract + content unit/integration + E2E |
| `getContestantAnnouncement` | `GET /api/v1/contestant/announcements/{announcement_slug}` | strict handler ✓ | service + store/adapter ✓ | Contestant /announces | contract + content unit/integration + E2E |
| `getContestantRanking` | `GET /api/v1/contestant/ranking` | strict handler ✓ | service + store/adapter ✓ | Contestant /ranking | contract + service/PostgreSQL + E2E |
| `getContestantRule` | `GET /api/v1/contestant/rule` | strict handler ✓ | service + store/adapter ✓ | Contestant /rule | contract + service/PostgreSQL + E2E |
| `getContestantDashboardSchedule` | `GET /api/v1/contestant/dashboard-schedule` | strict handler ✓ | service + store/adapter ✓ | Contestant dashboard | contract + service/PostgreSQL + E2E |
| `listAdminTeams` | `GET /api/v1/admin/teams` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `createAdminTeam` | `POST /api/v1/admin/teams` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `getAdminTeam` | `GET /api/v1/admin/teams/{team_code}` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `updateAdminTeam` | `PATCH /api/v1/admin/teams/{team_code}` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `deleteAdminTeam` | `DELETE /api/v1/admin/teams/{team_code}` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `listAdminInvitations` | `GET /api/v1/admin/invitations` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `createAdminInvitation` | `POST /api/v1/admin/invitations` | strict handler ✓ | service + store/adapter ✓ | Admin /teams | contract + service/PostgreSQL + E2E |
| `listAdminContestants` | `GET /api/v1/admin/contestants` | strict handler ✓ | service + store/adapter ✓ | Admin /contestants | contract + service/PostgreSQL + E2E |
| `getAdminContentStatus` | `GET /api/v1/admin/content/status` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `refreshAdminContent` | `POST /api/v1/admin/content/actions/refresh` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `listAdminSections` | `GET /api/v1/admin/sections` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `listAdminProblems` | `GET /api/v1/admin/problems` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `getAdminProblem` | `GET /api/v1/admin/problems/{problem_code}` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `listAdminAnnouncements` | `GET /api/v1/admin/announcements` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `getAdminAnnouncement` | `GET /api/v1/admin/announcements/{announcement_slug}` | strict handler ✓ | service + store/adapter ✓ | Admin /content | contract + content unit/integration + E2E |
| `listAdminAnswers` | `GET /api/v1/admin/answers` | strict handler ✓ | service + store/adapter ✓ | Admin /submissions | contract + service/PostgreSQL + E2E |
| `getAdminAnswer` | `GET /api/v1/admin/answers/{team_code}/{problem_code}/{answer_number}` | strict handler ✓ | service + store/adapter ✓ | Admin /submissions | contract + service/PostgreSQL + E2E |
| `listAdminMarkingResults` | `GET /api/v1/admin/marking-results` | strict handler ✓ | service + store/adapter ✓ | Admin /submissions | contract + service/PostgreSQL + E2E |
| `createAdminMarkingResult` | `POST /api/v1/admin/marking-results` | strict handler ✓ | service + store/adapter ✓ | Admin /submissions | contract + service/PostgreSQL + E2E |
| `listAdminScores` | `GET /api/v1/admin/scores` | strict handler ✓ | service + store/adapter ✓ | Admin /scores | contract + service/PostgreSQL + E2E |
| `recalculateAdminScores` | `POST /api/v1/admin/scores/actions/recalculate` | strict handler ✓ | service + store/adapter ✓ | Admin /scores | contract + service/PostgreSQL + E2E |
| `revealAdminFinalScores` | `POST /api/v1/admin/scores/actions/reveal-final` | strict handler ✓ | service + store/adapter ✓ | Admin /scores | contract + service/PostgreSQL + E2E |
| `listAdminDeployments` | `GET /api/v1/admin/deployments` | strict handler ✓ | service + store/adapter ✓ | Admin /deployments | contract + transition/SSE/Redis + E2E |
| `createAdminDeployment` | `POST /api/v1/admin/deployments` | strict handler ✓ | service + store/adapter ✓ | Admin /deployments | contract + transition/SSE/Redis + E2E |
| `streamAdminDeployments` | `GET /api/v1/admin/deployments/stream` | strict handler ✓ | service + store/adapter ✓ | Admin /deployments | contract + transition/SSE/Redis + E2E |
| `syncAdminDeployment` | `POST /api/v1/admin/deployments/{team_code}/{problem_code}/sync` | strict handler ✓ | service + store/adapter ✓ | Admin /deployments | contract + transition/SSE/Redis + E2E |
| `createAdminDeploymentEvent` | `POST /api/v1/admin/deployments/{team_code}/{problem_code}/{revision}/events` | strict handler ✓ | service + store/adapter ✓ | SState callback (server-to-server) | contract + transition/SSE/Redis + E2E |
| `getAdminRanking` | `GET /api/v1/admin/ranking` | strict handler ✓ | service + store/adapter ✓ | Admin /scores | contract + service/PostgreSQL + E2E |
| `getAdminRule` | `GET /api/v1/admin/rule` | strict handler ✓ | service + store/adapter ✓ | Admin /settings | contract + service/PostgreSQL + E2E |
| `replaceAdminRule` | `PUT /api/v1/admin/rule` | strict handler ✓ | service + store/adapter ✓ | Admin /settings | contract + service/PostgreSQL + E2E |
| `getAdminDashboardSchedule` | `GET /api/v1/admin/dashboard-schedule` | strict handler ✓ | service + store/adapter ✓ | Admin /settings | contract + service/PostgreSQL + E2E |
| `replaceAdminDashboardSchedule` | `PUT /api/v1/admin/dashboard-schedule` | strict handler ✓ | service + store/adapter ✓ | Admin /settings | contract + service/PostgreSQL + E2E |

## Cross-cutting acceptance gates

- Generated artifacts: `task generate` 後の `backend/internal/transport/api/api.gen.go` と `frontend/packages/api/src/schema.d.ts` に差分がないこと。
- Request/response contract: strict Chi handler、request validator、response validator、全operation登録テストを通すこと。
- Persistence: PostgreSQL/Redis integration testsはTestcontainersで分離実行し、競合・履歴・snapshot・Pub/Sub復元を検証すること。
- Redeploy: background pollingは禁止し、POST 202相当の受付、Bearer callback、event ID冪等化、SSE snapshot/event/reconnect、運営の一回syncだけを許可すること。
- UI: Adminは8業務route、Contestantは認証状態・profile・schedule・problem/answer・announcement・team・ranking・rule・activityをREST typed client経由で利用すること。
- Errors: RFC 9457の401/403/404/409/422/429/502を共通表示し、429はRetry-After countdown、mutation失敗時はformを保持すること。
