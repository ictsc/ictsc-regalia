import { api, expectData, expectNoContent, type ApiClient } from "@ictsc/api";

export type ImpersonationCandidate = {
  name: string;
  displayName: string;
  teamName: string;
  teamCode: number;
};

function candidateRequestError(action: string): Error {
  return new Error(`成り代わり${action}に失敗しました`);
}

export async function listImpersonationCandidates(
  client: ApiClient,
): Promise<ImpersonationCandidate[]> {
  const response = await client.GET("/api/v1/admin/contestants");
  const { contestants } = expectData(response);
  try {
    return contestants.map((contestant) => ({
      name: contestant.profile.name,
      displayName: contestant.profile.display_name,
      teamName: contestant.team.name,
      teamCode: Number(contestant.team.code),
    }));
  } catch (err) {
    if (err instanceof Error) throw err;
    throw candidateRequestError("候補の取得");
  }
}

export async function startImpersonation(candidate: {
  name: string;
  teamCode: number;
}): Promise<void> {
  expectNoContent(
    await api.POST("/api/v1/admin/impersonations", {
      body: { contestant_name: candidate.name },
    }),
  );
}
