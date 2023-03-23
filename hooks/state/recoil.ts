import { atom, useRecoilValue, useSetRecoilState } from "recoil";
import { setCookie, parseCookies } from "nookies";
import { useEffect } from "react";

export const dismissNoticeIdsState = atom<string[]>({
  key: "dismiss-notice-id",
  default: [],
});

export function DismissNoticeStateInit() {
  const setDismissNoticeIds = useSetRecoilState(dismissNoticeIdsState);

  useEffect(() => {
    //dismissNoticeIds を cookie から取得
    const cookies = parseCookies();

    if (cookies.dismissNoticeIds !== undefined) {
      const ids: string[] = JSON.parse(cookies.dismissNoticeIds);

      if (ids.length > 0) {
        setDismissNoticeIds(ids);
      }
    }
  }, [setDismissNoticeIds]);

  return null;
}

export function WatchDismissNoticeIds() {
  const dismissNoticeIds = useRecoilValue(dismissNoticeIdsState);

  useEffect(() => {
    setCookie(null, "dismissNoticeIds", JSON.stringify(dismissNoticeIds), {
      maxAge: 30 * 24 * 60 * 60 * 1000,
    });
  }, [dismissNoticeIds]);

  return null;
}
