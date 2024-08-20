import type versions from "@material-symbols/metadata/versions.json";
import clsx from "clsx";
import { CSSProperties } from "react";

export type IconTypes = keyof typeof versions;

export type MaterilaSymbolProps = {
  readonly icon: IconTypes;
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
    // material symbols の提供するクラスを使わないとフォントを指定できない
    // eslint-disable-next-line tailwindcss/no-custom-classname
    <span className={clsx("material-symbols-outlined select-none", className)} style={style}>
      {icon}
    </span>
  );
}
