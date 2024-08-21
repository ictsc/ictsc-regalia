import type { MaterialSymbol } from "material-symbols";
import clsx from "clsx";
import { CSSProperties } from "react";

export type { MaterialSymbol as MaterialSymbolType };

export type MaterilaSymbolProps = {
  readonly icon: MaterialSymbol;
  readonly fill?: boolean;
  readonly size?: number;
  readonly className?: string;
  readonly style?: CSSProperties;
};

export function MaterialSymbol({
  icon,
  fill,
  size,
  className,
  style: propStyle,
}: MaterilaSymbolProps) {
  const fontVariationSettings = [
    propStyle?.fontVariationSettings,
    size != null && `"opsz" ${size}`,
    fill && '"FILL" 1',
  ]
    .filter(Boolean)
    .join(",");
  const style: CSSProperties = {
    ...propStyle,
    fontVariationSettings,
    ...(size != null ? { fontSize: size } : {}),
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
