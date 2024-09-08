export type formType = "normal" | "multiple";

export interface Problem {
  id: string;
  code: string;
  title: string;
  body?: string;
  type: formType;
  point: number;
  solved_criterion: number | null;
  previous_problem_id: string | null;
  unchecked: number | null;
  unchecked_near_overdue: number | null;
  unchecked_overdue: number | null;
  current_point: number;
  is_answered: boolean;
  is_solved: boolean;
}

interface ConnectionInfo {
  type?: string;
  hostname?: string;
  port?: number;
  user?: string;
  password?: string;
  command?: string;
}

export interface Matter {
  code: string;
  title: string;
  point: number;
  solvedCriterion: number;
  authorId: string;
  connectInfo?: ConnectionInfo[];
}

export const testProblem: Problem = {
  id: "id",
  code: "XYZ",
  title: "テスト問題タイトル",
  body: "---\ntitle: テスト\n---\n# テスト本文",
  type: "normal",
  point: 100,
  solved_criterion: 150,
  previous_problem_id: null,
  unchecked: null,
  unchecked_near_overdue: null,
  unchecked_overdue: null,
  current_point: 100,
  is_answered: false,
  is_solved: false,
};
