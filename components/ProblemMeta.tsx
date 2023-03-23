import { Problem } from "@/types/Problem";

interface Props {
  problem: Problem;
}

const ProblemMeta = ({ problem }: Props) => {
  return (
    <div className={"flex flex-row items-end py-4"}>
      満点
      <span className={"sm:text-2xl"}> {problem.point} </span>pt 採点基準
      <span className={"sm:text-2xl"}> {problem.solved_criterion} </span>
      pt 問題コード
      <span className={"sm:text-2xl"}> {problem.code}</span>
    </div>
  );
};

export default ProblemMeta;
