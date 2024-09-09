import { ConnectionInfo } from "@/proto/admin/v1/problem_pb";

export type Props = {
  connectionInfos: ConnectionInfo[];
};

function ProblemConnectionInfo({ connectionInfos }: Props) {
  return (
    <div className="overflow-x-auto">
      <table className="table table-compact w-full">
        <thead>
          <tr>
            <th>ホスト名</th>
            <th>コマンド</th>
            <th>パスワード</th>
            <th>種類</th>
          </tr>
        </thead>
        <tbody>
          {connectionInfos?.map((info) => (
            <tr key={info.hostname}>
              <th>{info.hostname}</th>
              <th>{info.command}</th>
              <th>{info.password}</th>
              <th>{info.type}</th>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ProblemConnectionInfo;
