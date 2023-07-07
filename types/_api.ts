import { Answer } from "@/types/Answer";
import { Problem } from "@/types/Problem";
import { Rank } from "@/types/Rank";
import { User } from "@/types/User";

export type Result<T> = {
  code: number;
  data: T | null;
};

// /api/auth/self
export type AuthSelfResult = {
  user: User;
};

// /api/problems
export type ProblemResult = {
  problems: Problem[];
};

// /api/ranking
export type RankingResult = {
  ranking: Rank[];
};

// /api/problems/:id/answers
export type AnswerResult = {
  answers: Answer[];
};

// /api/auth/signup
export type SignUpRequest = {
  name: string;
  password: string;
  user_group_id: string;
  invitation_code: string;
};
