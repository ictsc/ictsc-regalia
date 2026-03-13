import type { Meta, StoryObj } from "@storybook/react";
import { create } from "@bufbuild/protobuf";
import { timestampFromDate } from "@bufbuild/protobuf/wkt";
import {
  ScheduleEntrySchema,
  SubmissionStatusSchema,
  type Problem,
  type ScheduleEntry,
} from "@ictsc/proto/contestant/v1";
import { ProblemsPage } from "./page";

export default {
  title: "pages/problems",
  component: ProblemsPage,
} satisfies Meta<typeof ProblemsPage>;

type Story = StoryObj<typeof ProblemsPage>;

const day1Am = create(ScheduleEntrySchema, {
  name: "day1-am",
  startAt: timestampFromDate(new Date("2026-01-01T10:00:00+09:00")),
  endAt: timestampFromDate(new Date("2026-01-01T12:00:00+09:00")),
});

const day1Pm = create(ScheduleEntrySchema, {
  name: "day1-pm",
  startAt: timestampFromDate(new Date("2026-01-01T13:00:00+09:00")),
  endAt: timestampFromDate(new Date("2099-12-31T23:59:59+09:00")),
});

const day2Am = create(ScheduleEntrySchema, {
  name: "day2-am",
  startAt: timestampFromDate(new Date("2100-01-01T10:00:00+09:00")),
  endAt: timestampFromDate(new Date("2100-01-01T12:00:00+09:00")),
});

const day2Pm = create(ScheduleEntrySchema, {
  name: "day2-pm",
  startAt: timestampFromDate(new Date("2100-01-01T13:00:00+09:00")),
  endAt: timestampFromDate(new Date("2100-01-01T16:00:00+09:00")),
});

const submittable = create(SubmissionStatusSchema, {
  isSubmittable: true,
});

const notSubmittable = create(SubmissionStatusSchema, {
  isSubmittable: false,
});

function makeProblem(
  code: string,
  title: string,
  opts?: {
    score?: { score: number; markedScore: number; penalty: number };
    submissionStatus?: Problem["submissionStatus"];
    submissionableSchedules?: ScheduleEntry[];
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
    submissionableSchedules: opts?.submissionableSchedules ?? [],
  } as Problem;
}

export const Default: Story = {
  args: {
    notices: [],
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

export const GroupedBySchedule: Story = {
  args: {
    notices: [],
    problems: [
      // day1-am（過去 → 提出不可）
      makeProblem("NET01", "ネットワーク基礎問題", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day1Am],
        score: { score: 200, markedScore: 200, penalty: 0 },
      }),
      makeProblem("NET02", "VLAN設定問題", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day1Am],
        score: { score: 100, markedScore: 120, penalty: -20 },
      }),
      makeProblem("SRV01", "Webサーバー構築", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day1Am],
      }),
      // day1-pm（現在 → 提出可能）
      makeProblem("SRV02", "データベース復旧", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
      }),
      makeProblem("SRV03", "コンテナ運用管理", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
        score: { score: 80, markedScore: 100, penalty: -20 },
      }),
      makeProblem("DNS01", "DNS権威サーバー設定", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
      }),
      makeProblem("DNS02", "DNSキャッシュ問題", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
      }),
      // day2-am（未来 → 提出不可）
      makeProblem("SEC01", "セキュリティ診断", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day2Am],
      }),
      makeProblem("SEC02", "ファイアウォール設定", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day2Am],
      }),
      // day2-pm（未来 → 提出不可）
      makeProblem("APP01", "ロードバランサ冗長化", {
        submissionStatus: notSubmittable,
        submissionableSchedules: [day2Pm],
      }),
      // day1-pm + day2-am（複数スケジュール）
      makeProblem("MON01", "監視システム構築", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm, day2Am],
      }),
      // 全スケジュール
      makeProblem("ALL01", "総合演習", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Am, day1Pm, day2Am, day2Pm],
      }),
    ],
  },
};

export const SingleSchedule: Story = {
  args: {
    notices: [],
    problems: [
      makeProblem("NET01", "ネットワーク基礎問題", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
      }),
      makeProblem("NET02", "VLAN設定問題", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
        score: { score: 100, markedScore: 120, penalty: -20 },
      }),
      makeProblem("SRV01", "Webサーバー構築", {
        submissionStatus: submittable,
        submissionableSchedules: [day1Pm],
      }),
    ],
  },
};

export const NoSchedules: Story = {
  args: {
    notices: [],
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
