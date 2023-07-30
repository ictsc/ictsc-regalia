interface Props {
  className?: "" | "text-warning" | "text-error";
  value: number;
}

function ConditionalText({ className, value }: Props) {
  if (value > 0) {
    return <div className={`inline-block ${className}`}>{value}</div>;
  }

  return <div className="inline-block">-</div>;
}

ConditionalText.defaultProps = {
  className: "",
};

export default ConditionalText;
