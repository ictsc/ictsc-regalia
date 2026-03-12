import type { Meta, StoryObj } from "@storybook/react";
import { create } from "@bufbuild/protobuf";
import {
  SubmissionStatusSchema,
  type Problem,
} from "@ictsc/proto/contestant/v1";
import { ProblemsPage } from "./page";

export default {
  title: "pages/problems",
  component: ProblemsPage,
} satisfies Meta<typeof ProblemsPage>;

type Story = StoryObj<typeof ProblemsPage>;

function makeProblem(
  code: string,
  title: string,
  opts?: {
    score?: { score: number; markedScore: number; penalty: number };
    submissionStatus?: Problem["submissionStatus"];
  },
): Problem {
  return {
    $typeName: "contestant.v1.Problem" as const,
    code,
    title,
    maxScore: 200,
    category: "Network",
    score: opts?.score
      ? {
          $typeName: "contestant.v1.Score" as const,
          maxScore: 200,
          ...opts.score,
        }
      : undefined,
    submissionStatus: opts?.submissionStatus,
  } as Problem;
}

const submittable = create(SubmissionStatusSchema, {
  isSubmittable: true,
});

const notSubmittable = create(SubmissionStatusSchema, {
  isSubmittable: false,
});

export const Default: Story = {
  args: {
    problems: [
      ...Array.from({ length: 4 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Network",
        submissionableSchedules: [],
      })),
      ...Array.from({ length: 5 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Server",
        submissionableSchedules: [],
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 100,
          markedScore: 120,
          penalty: -20,
        },
      })),
      ...Array.from({ length: 2 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Network",
        submissionableSchedules: [],
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 160,
          markedScore: 200,
          penalty: -40,
        },
      })),
      ...Array.from({ length: 5 }, () => ({
        $typeName: "contestant.v1.Problem" as const,
        code: "ABC",
        title: "問題 ABC",
        maxScore: 200,
        category: "Server",
        submissionableSchedules: [],
        score: {
          $typeName: "contestant.v1.Score" as const,
          score: 200,
          markedScore: 200,
          penalty: 0,
        },
      })),
    ],
  },
};

export const Grouped: Story = {
  args: {
    problems: [
      // 提出可能
      makeProblem("NET01", "ネットワーク基礎問題", {
        submissionStatus: submittable,
      }),
      makeProblem("NET02", "VLAN設定問題", {
        submissionStatus: submittable,
        score: { score: 100, markedScore: 120, penalty: -20 },
      }),
      makeProblem("SRV01", "Webサーバー構築", {
        submissionStatus: submittable,
      }),
      makeProblem("SRV02", "データベース復旧", {
        submissionStatus: submittable,
      }),
      makeProblem("SRV03", "コンテナ運用管理", {
        submissionStatus: submittable,
        score: { score: 80, markedScore: 100, penalty: -20 },
      }),
      makeProblem("DNS01", "DNS権威サーバー設定", {
        submissionStatus: submittable,
      }),
      makeProblem("DNS02", "DNSキャッシュ問題", {
        submissionStatus: submittable,
      }),
      makeProblem("MON01", "監視システム構築", {
        submissionStatus: submittable,
      }),
      // 提出不可
      makeProblem("SEC01", "セキュリティ診断", {
        submissionStatus: notSubmittable,
      }),
      makeProblem("SEC02", "ファイアウォール設定", {
        submissionStatus: notSubmittable,
      }),
      makeProblem("SEC03", "IDS/IPS チューニング", {
        submissionStatus: notSubmittable,
      }),
      makeProblem("APP01", "ロードバランサ冗長化", {
        submissionStatus: notSubmittable,
      }),
      makeProblem("OLD01", "終了済み：OSPF経路制御", {
        submissionStatus: notSubmittable,
        score: { score: 200, markedScore: 200, penalty: 0 },
      }),
      makeProblem("OLD02", "終了済み：BGPピアリング", {
        submissionStatus: notSubmittable,
        score: { score: 150, markedScore: 180, penalty: -30 },
      }),
      makeProblem("OLD03", "終了済み：IPv6移行", {
        submissionStatus: notSubmittable,
        score: { score: 60, markedScore: 80, penalty: -20 },
      }),
      makeProblem("OLD04", "終了済み：RADIUS認証", {
        submissionStatus: notSubmittable,
      }),
    ],
  },
};

export const SubmittableOnly: Story = {
  args: {
    problems: [
      makeProblem("NET01", "ネットワーク基礎問題", {
        submissionStatus: submittable,
      }),
      makeProblem("NET02", "VLAN設定問題", {
        submissionStatus: submittable,
        score: { score: 100, markedScore: 120, penalty: -20 },
      }),
      makeProblem("SRV01", "Webサーバー構築", {
        submissionStatus: submittable,
      }),
    ],
  },
};

export const NotSubmittableOnly: Story = {
  args: {
    problems: [
      makeProblem("OLD01", "終了済み：OSPF経路制御", {
        submissionStatus: notSubmittable,
        score: { score: 200, markedScore: 200, penalty: 0 },
      }),
      makeProblem("OLD02", "終了済み：BGPピアリング", {
        submissionStatus: notSubmittable,
        score: { score: 150, markedScore: 180, penalty: -30 },
      }),
    ],
  },
};
