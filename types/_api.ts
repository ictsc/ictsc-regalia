import { Answer } from "@/types/Answer";
import { Problem } from "@/types/Problem";
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

// /api/problems/:id/answers
export type AnswerResult = {
  answer: Answer;
};
