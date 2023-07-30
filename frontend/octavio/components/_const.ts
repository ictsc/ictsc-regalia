const replaceN = (str: string) => str.replace(/\\n/g, "\n");

export const apiUrl = process.env.NEXT_PUBLIC_API_URL;
export const site = process.env.NEXT_PUBLIC_SITE_NAME;
export const rule = replaceN(process.env.RULE ?? "");
export const shortRule = replaceN(process.env.NEXT_PUBLIC_SHORT_RULE ?? "");
export const recreateRule = replaceN(
  process.env.NEXT_PUBLIC_RECREATE_RULE ?? ""
);
export const answerLimit = process.env.NEXT_PUBLIC_ANSWER_LIMIT;
