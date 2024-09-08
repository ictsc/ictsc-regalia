"use client";

import { useParams } from "next/navigation";

import { useMutation, useQuery } from "@connectrpc/connect-query";

import DescriptiveProblemEdit from "@/app/problems/__components/DescriptiveProblemEdit";
import {
  getProblem,
  patchProblem,
} from "@/proto/admin/v1/problem-ProblemService_connectquery";

function Index() {
  const { id } = useParams<{ id: string }>();

  const result = useQuery(getProblem, { id });
  const mutation = useMutation(patchProblem);

  if (result.data?.problem?.body.case !== "descriptive") {
    return <div data-testid="unsupported">この問題タイプは未対応です。</div>;
  }

  return (
    <DescriptiveProblemEdit
      problem={result.data?.problem}
      saveProblemData={async (data) =>
        mutation
          .mutateAsync({
            id,
            title: data.title,
            code: data.code,
            point: data.point,
            body: {
              case: "descriptive",
              value: {
                body: data.body,
                connectionInfos: data.connectionInfos,
              },
            },
          })
          .then(() => true)
          .catch(() => false)
      }
    />
  );
}

export default Index;
