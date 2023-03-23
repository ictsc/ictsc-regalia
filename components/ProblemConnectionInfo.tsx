import HoverCopyText from "@/components/HoverCopyText";
import { Matter } from "@/types/Problem";

interface Props {
  matter: Matter | null;
}

function ProblemConnectionInfo({ matter }: Props) {
  return (
    <div className="overflow-x-auto">
      <table className="table table-compact w-full">
        <thead>
          <tr>
            <th>ホスト名</th>
            <th>コマンド</th>
            <th>ユーザ</th>
            <th>パスワード</th>
            <th>ポート</th>
            <th>種類</th>
          </tr>
        </thead>
        <tbody>
          {matter?.connectInfo?.map((info, index) => (
            // eslint-disable-next-line react/no-array-index-key
            <tr key={index}>
              <th>{info.hostname}</th>
              <HoverCopyText text={info.command ?? ""} />
              <HoverCopyText text={info.user ?? ""} />
              <HoverCopyText text={info.password ?? ""} />
              <td>{info.port}</td>
              <td>{info.type}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default ProblemConnectionInfo;
