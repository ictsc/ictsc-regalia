import { startTransition, use, useDeferredValue, useState } from "react";
import {
  createFileRoute,
  useNavigate,
  useRouter,
} from "@tanstack/react-router";
import {
  listImpersonationCandidates,
  startImpersonation,
  type ImpersonationCandidate,
} from "@app/features/viewer";
import { ImpersonationPage } from "./impersonation.page";

export const Route = createFileRoute("/signin/impersonation")({
  component: RouteComponent,
  loader: () => ({
    candidates: listImpersonationCandidates(),
  }),
});

function RouteComponent() {
  const { candidates: candidatesPromise } = Route.useLoaderData();
  const candidates = use(useDeferredValue(candidatesPromise));
  const [selectedKey, setSelectedKey] = useState(() =>
    candidateKey(candidates[0]),
  );
  const [error, setError] = useState<string>();
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const router = useRouter();

  async function handleStartImpersonation() {
    const selected = candidates.find((c) => candidateKey(c) === selectedKey);
    if (selected == null) return;

    setLoading(true);
    setError(undefined);
    try {
      await startImpersonation({
        name: selected.name,
        teamCode: selected.teamCode,
      });
      await router.invalidate();
      startTransition(() => navigate({ to: "/" }));
    } catch (err) {
      setError(err instanceof Error ? err.message : "偽装の開始に失敗しました");
    } finally {
      setLoading(false);
    }
  }

  return (
    <ImpersonationPage>
      {candidates.length === 0 ? (
        <p className="text-14 text-text">偽装できる競技者が見つかりません。</p>
      ) : (
        <div className="flex flex-col gap-12">
          {error != null ? (
            <p className="text-14 text-primary">{error}</p>
          ) : null}
          <label className="flex flex-col gap-4">
            <span className="text-14 text-text font-bold">偽装する競技者</span>
            <select
              value={selectedKey}
              onChange={(event) => setSelectedKey(event.target.value)}
              className="border-disabled rounded-4 text-14 text-text w-full border px-12 py-8"
            >
              {candidates.map((candidate) => (
                <option
                  key={`${candidate.teamCode}:${candidate.name}`}
                  value={`${candidate.teamCode}:${candidate.name}`}
                >
                  {candidate.teamCode} {candidate.teamName} /{" "}
                  {candidate.displayName || candidate.name}
                </option>
              ))}
            </select>
          </label>
          <button
            type="button"
            disabled={selectedKey === "" || loading}
            onClick={() => {
              void handleStartImpersonation();
            }}
            className="bg-primary text-surface-0 rounded-8 text-14 px-16 py-8 hover:opacity-80 disabled:opacity-40"
          >
            この競技者として入る
          </button>
        </div>
      )}
    </ImpersonationPage>
  );
}

function candidateKey(
  candidate?: Pick<ImpersonationCandidate, "name" | "teamCode">,
) {
  if (candidate == null) return "";
  return `${candidate.teamCode}:${candidate.name}`;
}
