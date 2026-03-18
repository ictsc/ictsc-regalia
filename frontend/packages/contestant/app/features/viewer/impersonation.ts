import {
  type Transport,
  createClient,
  ConnectError,
} from "@connectrpc/connect";
import { ViewerService } from "@ictsc/proto/contestant/v1";

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
  transport: Transport,
): Promise<ImpersonationCandidate[]> {
  const client = createClient(ViewerService, transport);
  try {
    const { contestants } = await client.listContestants({});
    return contestants.map((contestant) => ({
      name: contestant.name,
      displayName: contestant.displayName,
      teamName: contestant.teamName,
      teamCode: Number(contestant.teamCode),
    }));
  } catch (err) {
    if (err instanceof ConnectError) throw err;
    throw candidateRequestError("候補の取得");
  }
}

export async function startImpersonation(candidate: {
  name: string;
  teamCode: number;
}): Promise<void> {
  const response = await fetch("/api/auth/impersonation", {
    method: "POST",
    headers: {
      "content-type": "application/json",
    },
    body: JSON.stringify(candidate),
  });
  if (!response.ok) {
    throw candidateRequestError("開始");
  }
}
