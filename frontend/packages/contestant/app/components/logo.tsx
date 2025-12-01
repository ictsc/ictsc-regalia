import pathToLogo from "@assets/ictsc2025.png";

const RATIO = 1000 / 267;

export type LogoProps = {
  readonly className?: string;
  readonly width?: number;
  readonly height?: number;
};

export function Logo({ className, width, height }: LogoProps) {
  if (width != null && height == null) {
    height = width / RATIO;
  }
  if (height != null && width == null) {
    width = height * RATIO;
  }
  return (
    <img
      className={className}
      width={width}
      height={height}
      src={pathToLogo}
      alt="ICTSC: ICT Trouble Shooting Contest"
    />
  );
}
