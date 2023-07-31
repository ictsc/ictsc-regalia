"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";

import { mutate } from "swr";

import useAuth from "@/hooks/auth";

function ICTSCNavBar() {
  const router = useRouter();

  const { user, logout } = useAuth();

  const handleLogout = async () => {
    const response = await logout();

    if (response.code === 200) {
      await mutate(() => true, undefined, { revalidate: true });
      await router.push("/");
    }
  };

  return (
    <div className="navbar bg-primary text-primary-content">
      <div className="flex-1">
        <Link href="/" className="btn btn-ghost normal-case text-xl">
          ICTSC
        </Link>
      </div>
      <nav className="flex-none">
        <ul className="menu menu-horizontal p-0">
          <li>
            <Link href="/">ルール</Link>
          </li>
          {user !== null && (
            <>
              <li>
                <Link href="/team_info">チーム情報</Link>
              </li>
              <li>
                <Link href="/problems">問題</Link>
              </li>
            </>
          )}
          <li>
            <Link href="/ranking">順位</Link>
          </li>
          {user !== null && (
            <>
              <li>
                <Link href="/users">参加者</Link>
              </li>
              {user.user_group.is_full_access && !user.is_read_only && (
                <li>
                  <Link href="/scoring">採点</Link>
                </li>
              )}
            </>
          )}
          {user === null ? (
            <li className="ml-4">
              <Link href="/login">ログイン</Link>
            </li>
          ) : (
            // eslint-disable-next-line jsx-a11y/no-noninteractive-tabindex
            <li tabIndex={0} className="ml-4 dropdown dropdown-end">
              <div>{user.display_name}</div>
              <ul
                /* eslint-disable-next-line jsx-a11y/no-noninteractive-tabindex */
                tabIndex={0}
                className="menu menu-compact dropdown-content bg-base-100 mt-3 p-2 shadow rounded-box w-52 text-base-content"
              >
                <li>
                  <div>チーム: {user.user_group.name}</div>
                </li>
                <li>
                  <Link href="/profile">プロフィール</Link>
                </li>
                <div className="divider my-0" />
                <li>
                  <button type="button" onClick={handleLogout}>
                    ログアウト
                  </button>
                </li>
              </ul>
            </li>
          )}
        </ul>
      </nav>
    </div>
  );
}

export default ICTSCNavBar;
