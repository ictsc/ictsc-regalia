import { minLength, object, Output, string, StringSchema } from "valibot";

const nameSchema: StringSchema = string([
  minLength(1, "ユーザー名を入力してください"),
]);

const passwordSchema: StringSchema = string([
  minLength(1, "パスワードを入力してください"),
]);

export const LoginSchema = object({
  name: nameSchema,
  password: passwordSchema,
});

export type LoginType = Output<typeof LoginSchema>;
