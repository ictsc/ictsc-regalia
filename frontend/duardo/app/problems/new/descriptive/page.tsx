"use client";

import { useMutation } from "@connectrpc/connect-query";

import DescriptiveProblemEdit from "@/app/problems/__components/DescriptiveProblemEdit";
import { postProblem } from "@/proto/admin/v1/problem-ProblemService_connectquery";

function Index() {
  const mutation = useMutation(postProblem);

  return (
    <DescriptiveProblemEdit
      saveProblemData={async (data) =>
        mutation
          .mutateAsync({
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
