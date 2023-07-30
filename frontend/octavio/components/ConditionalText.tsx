function ConditionalText({
  className = "",
  value,
}: {
  // eslint-disable-next-line react/require-default-props
  className?: "" | "text-warning" | "text-error";
  value: number;
}) {
  if (value > 0) {
    return <div className={`inline-block ${className}`}>{value}</div>;
  }

  return <div className="inline-block">-</div>;
}

export default ConditionalText;
