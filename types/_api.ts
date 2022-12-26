import {User} from "./User";
import {Problem} from "./Problem";

export type Result<T> = {
  code: number;
  data: T | null;
}

// /api/auth/self
export type AuthSelfResult = {
  user: User;
}

// /api/problems
export type ProblemResult = {
  problems: Problem[];
}