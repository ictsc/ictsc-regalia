import Link from "next/link";
import { useRouter } from "next/router";

import { useApi } from "../hooks/api";
import { useAuth } from "../hooks/auth";

const ICTSCNavBar = () => {
  const router = useRouter();

  const { apiClient } = useApi();
  const { user, mutate } = useAuth();

  const logout = async () => {
    const response = await apiClient.delete("auth/signout");

    if (response.ok) {
      await mutate();
      await router.push("/");
    }
  };

  return (
    <div className={"navbar bg-primary text-primary-content"}>
      <div className={"flex-1"}>
        <Link href={"/"} className="btn btn-ghost normal-case text-xl">
          ICTSC
        </Link>
      </div>
      <div className={"flex-none"}>
        <ul className="menu menu-horizontal p-0">
          <li>
            <Link href={"/"}>ルール</Link>
          </li>
          {user !== null && (
            <>
              <li>
                <Link href={"/team_info"}>チーム情報</Link>
              </li>
              <li>
                <Link href={"/problems"}>問題</Link>
              </li>
              <li>
                <Link href={"/ranking"}>順位</Link>
              </li>
              <li>
                <Link href={"/users"}>参加者</Link>
              </li>
              {user.user_group.is_full_access && (
                <li>
                  <Link href={"/scoring"}>採点</Link>
                </li>
              )}
            </>
          )}
          {user === null ? (
            <li className={"ml-4"}>
              <Link href={"/login"}>ログイン</Link>
            </li>
          ) : (
            <li tabIndex={0} className={"ml-4 dropdown dropdown-end"}>
              <div>{user.display_name}</div>
              <ul
                tabIndex={0}
                className="menu menu-compact dropdown-content mt-3 p-2 shadow rounded-box w-52 text-base-content"
              >
                <li>
                  <a>チーム: {user.user_group.name}</a>
                </li>
                <li>
                  <Link href={"/profile"}>プロフィール</Link>
                </li>
                <div className="divider my-0" />
                <li>
                  <a onClick={logout}>ログアウト</a>
                </li>
              </ul>
            </li>
          )}
        </ul>
      </div>
    </div>
  );
};

export default ICTSCNavBar;
