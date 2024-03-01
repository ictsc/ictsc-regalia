import { z } from "zod";

import { ConnectionInfo } from "@/proto/admin/v1/problem_pb";

export const formSchema = z.object({
  title: z.string({ required_error: "タイトルは必須です" }).min(1, {
    message: "タイトルは1文字以上でなければなりません",
  }),
  code: z
    .string()
    .refine(
      (code) => /^[A-Z]{3}$/.test(code),
      "問題コードは大文字A-Zの3文字でなければなりません",
    ),
  point: z
    .number({ invalid_type_error: "ポイントは必須です" })
    .min(0, { message: "ポイントは0以上でなければなりません" })
    .refine(
      (point) => Number.isInteger(point),
      "ポイントは整数でなければなりません",
    ),
  connectionInfos: z
    .array(
      z.object({
        hostname: z.string().nullable(),
        command: z.string().nullable(),
        password: z.string().nullable(),
        type: z.string().nullable(),
      }),
    )
    .nullable()
    .transform(
      (infos): ConnectionInfo[] =>
        infos?.map((info) => {
          const connectionInfo = new ConnectionInfo();
          connectionInfo.hostname = info.hostname ?? "";
          connectionInfo.command = info.command ?? "";
          connectionInfo.password = info.password ?? "";
          connectionInfo.type = info.type ?? "";
          return connectionInfo;
        }) ?? [],
    ),
  body: z.string(),
});

export type FormSchema = z.infer<typeof formSchema>;
