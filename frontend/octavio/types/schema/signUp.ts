import { minLength, object, Output, string, StringSchema } from "valibot";

const nameSchema: StringSchema = string([
  minLength(1, "ユーザー名を入力してください"),
]);

const passwordSchema: StringSchema = string([
  minLength(8, "パスワードは8文字以上である必要があります"),
]);

export const SignUpSchema = object({
  name: nameSchema,
  password: passwordSchema,
});

export type SignUpType = Output<typeof SignUpSchema>;
