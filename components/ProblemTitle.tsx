interface Props {
  title: string;
}

const ProblemTitle = ({ title }: Props) => {
  return <h1 className={"title-ictsc pr-2 text-3xl"}>{title}</h1>;
};

export default ProblemTitle;
