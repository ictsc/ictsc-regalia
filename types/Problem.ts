export interface Problem {
  id: string;
  code: string;
  title: string;
  body?: string;
  point: number;
  solved_criterion: number | null;
  previous_problem_id: string | null;
  author_id: string;
  unchecked: number | null;
  unchecked_near_overdue: number | null;
  unchecked_overdue: number | null;
}

export interface Matter {
  code: string;
  title: string;
  point: number;
  solvedCriterion: number;
  authorId: string;
  connectInfo?: ConnectionInfo[];
}


interface ConnectionInfo {
  type?: string;
  hostname?: string;
  port?: number;
  user?: string;
  password?: string;
  command?: string;
}