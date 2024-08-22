import type { CSSProperties } from "react";
import type { MaterialSymbol } from "material-symbols";
import clsx from "clsx";

export type { MaterialSymbol as MaterialSymbolType };

export type MaterilaSymbolProps = {
  readonly icon: MaterialSymbol;
  readonly fill?: boolean;
  readonly size?: 20 | 24 | 40 | 48;
  readonly className?: string;
  readonly style?: CSSProperties;
};

export function MaterialSymbol({
  icon,
  fill = false,
  size = 24,
  className,
  style: propStyle,
}: MaterilaSymbolProps) {
  const style: CSSProperties = {
    ...propStyle,
    fontVariationSettings: `"FILL" ${fill ? "1" : "0"}`,
    fontSize: size,
  };
  return (
    <span
      // material symbols の提供するクラスを使わないとフォントを指定できない
      // eslint-disable-next-line tailwindcss/no-custom-classname
      className={clsx("material-symbols-outlined select-none", className)}
      style={style}
    >
      {icon}
    </span>
  );
}
