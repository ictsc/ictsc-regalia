import Error from "next/error";

import ICTSCNavBar from "../components/Navbar";
import LoadingPage from "../components/LoadingPage";
import { useUserGroups } from "../hooks/userGroups";
import Head from "next/head";
import { site } from "../components/_const";

const Users = () => {
  const { userGroups, isLoading } = useUserGroups();

  if (isLoading) {
    return (
      <>
        <ICTSCNavBar />
        <LoadingPage />
      </>
    );
  }

  if (userGroups === null) {
    return <Error statusCode={404} />;
  }

  return (
    <>
      <Head>
        <title>参加者一覧 - {site}</title>
      </Head>
      <ICTSCNavBar />
      <h1 className={"title-ictsc text-center py-12"}>参加者一覧</h1>
      <div className={"container-ictsc"}>
        <table className={"table border rounded-md w-full"}>
          <thead>
            <tr>
              <th>名前</th>
              <th>チーム名</th>
              <th>自己紹介</th>
            </tr>
          </thead>
          <tbody className={"text-sm"}>
            {userGroups?.map((userGroup) =>
              userGroup.members?.map((member) => (
                <tr key={member.id}>
                  <td className={"whitespace-normal max-w-[300px]"}>
                    {member.display_name}
                    <div className={"flex flex-row"}>
                      {(member.profile?.github_id ?? "") != "" && (
                        <a
                          href={`https://github.com/${member.profile?.github_id}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className={"btn btn-circle btn-ghost btn-xs"}
                        >
                          <GithubIcon />
                        </a>
                      )}
                      {(member.profile?.twitter_id ?? "") != "" && (
                        <a
                          href={`https://twitter.com/${member.profile?.twitter_id}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className={"btn btn-circle btn-ghost btn-xs"}
                        >
                          <TwitterIcon />
                        </a>
                      )}
                      {(member.profile?.facebook_id ?? "") != "" && (
                        <a
                          href={`https://www.facebook.com/${member.profile?.facebook_id}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className={"btn btn-circle btn-ghost btn-xs"}
                        >
                          <FacebookIcon />
                        </a>
                      )}
                    </div>
                  </td>
                  <td className={"whitespace-normal lg:min-w-[196px]"}>
                    {userGroup.name}
                  </td>
                  <td className={"whitespace-normal"}>
                    {member.profile?.self_introduction}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </>
  );
};

const GithubIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    height="20px"
    width="20px"
  >
    <g data-name="Layer 2">
      <rect width="24" height="24" opacity="0" />
      <path d="M16.24 22a1 1 0 0 1-1-1v-2.6a2.15 2.15 0 0 0-.54-1.66 1 1 0 0 1 .61-1.67C17.75 14.78 20 14 20 9.77a4 4 0 0 0-.67-2.22 2.75 2.75 0 0 1-.41-2.06 3.71 3.71 0 0 0 0-1.41 7.65 7.65 0 0 0-2.09 1.09 1 1 0 0 1-.84.15 10.15 10.15 0 0 0-5.52 0 1 1 0 0 1-.84-.15 7.4 7.4 0 0 0-2.11-1.09 3.52 3.52 0 0 0 0 1.41 2.84 2.84 0 0 1-.43 2.08 4.07 4.07 0 0 0-.67 2.23c0 3.89 1.88 4.93 4.7 5.29a1 1 0 0 1 .82.66 1 1 0 0 1-.21 1 2.06 2.06 0 0 0-.55 1.56V21a1 1 0 0 1-2 0v-.57a6 6 0 0 1-5.27-2.09 3.9 3.9 0 0 0-1.16-.88 1 1 0 1 1 .5-1.94 4.93 4.93 0 0 1 2 1.36c1 1 2 1.88 3.9 1.52a3.89 3.89 0 0 1 .23-1.58c-2.06-.52-5-2-5-7a6 6 0 0 1 1-3.33.85.85 0 0 0 .13-.62 5.69 5.69 0 0 1 .33-3.21 1 1 0 0 1 .63-.57c.34-.1 1.56-.3 3.87 1.2a12.16 12.16 0 0 1 5.69 0c2.31-1.5 3.53-1.31 3.86-1.2a1 1 0 0 1 .63.57 5.71 5.71 0 0 1 .33 3.22.75.75 0 0 0 .11.57 6 6 0 0 1 1 3.34c0 5.07-2.92 6.54-5 7a4.28 4.28 0 0 1 .22 1.67V21a1 1 0 0 1-.94 1z" />
    </g>
  </svg>
);

const TwitterIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    height="20px"
    width="20px"
  >
    <g data-name="Layer 2">
      <g data-name="twitter">
        <polyline points="0 0 24 0 24 24 0 24" opacity="0" />
        <path d="M8.51 20h-.08a10.87 10.87 0 0 1-4.65-1.09A1.38 1.38 0 0 1 3 17.47a1.41 1.41 0 0 1 1.16-1.18 6.63 6.63 0 0 0 2.54-.89 9.49 9.49 0 0 1-3.51-9.07 1.41 1.41 0 0 1 1-1.15 1.35 1.35 0 0 1 1.43.41 7.09 7.09 0 0 0 4.88 2.75 4.5 4.5 0 0 1 1.41-3.1 4.47 4.47 0 0 1 6.37.19.7.7 0 0 0 .78.1A1.39 1.39 0 0 1 21 7.13a6.66 6.66 0 0 1-1.28 2.6A10.79 10.79 0 0 1 8.51 20zm0-2h.08a8.79 8.79 0 0 0 9.09-8.59 1.32 1.32 0 0 1 .37-.85 5.19 5.19 0 0 0 .62-1 2.56 2.56 0 0 1-1.91-.85A2.45 2.45 0 0 0 15 6a2.5 2.5 0 0 0-1.79.69 2.53 2.53 0 0 0-.72 2.42l.26 1.14-1.17.08a8.3 8.3 0 0 1-6.54-2.4 7.12 7.12 0 0 0 3.73 6.46l.95.54-.63.9a5.62 5.62 0 0 1-2.68 1.92A8.34 8.34 0 0 0 8.5 18zM19 6.65z" />
      </g>
    </g>
  </svg>
);

const FacebookIcon = () => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    height="20px"
    width="20px"
  >
    <g data-name="Layer 2">
      <g data-name="facebook">
        <rect
          width="24"
          height="24"
          transform="rotate(180 12 12)"
          opacity="0"
        />
        <rect
          width="24"
          height="24"
          transform="rotate(180 12 12)"
          opacity="0"
        />
        <path d="M13 22H9a1 1 0 0 1-1-1v-6.2H6a1 1 0 0 1-1-1v-3.6a1 1 0 0 1 1-1h2V7.5A5.77 5.77 0 0 1 14 2h3a1 1 0 0 1 1 1v3.6a1 1 0 0 1-1 1h-3v1.6h3a1 1 0 0 1 .8.39 1 1 0 0 1 .16.88l-1 3.6a1 1 0 0 1-1 .73H14V21a1 1 0 0 1-1 1zm-3-2h2v-6.2a1 1 0 0 1 1-1h2.24l.44-1.6H13a1 1 0 0 1-1-1V7.5a2 2 0 0 1 2-1.9h2V4h-2a3.78 3.78 0 0 0-4 3.5v2.7a1 1 0 0 1-1 1H7v1.6h2a1 1 0 0 1 1 1z" />
      </g>
    </g>
  </svg>
);

export default Users;
