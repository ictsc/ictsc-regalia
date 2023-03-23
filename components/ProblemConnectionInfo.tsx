import HoverCopyText from "@/components/HoverCopyText";
import { Matter } from "@/types/Problem";

interface Props {
  matter: Matter | null;
}

const ProblemConnectionInfo = ({ matter }: Props) => {
  return (
    <div className="overflow-x-auto">
      <table className="table table-compact w-full">
        <thead>
          <tr>
            <th>ホスト名</th>
            <th></th>
            <th>ユーザ</th>
            <th>パスワード</th>
            <th>ポート</th>
            <th>種類</th>
          </tr>
        </thead>
        <tbody>
          {matter?.connectInfo?.map((info, index) => (
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
};

export default ProblemConnectionInfo;
