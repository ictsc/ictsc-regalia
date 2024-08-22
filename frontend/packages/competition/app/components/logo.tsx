import pathToLogo from "../../assets/ictsc.png";

const RATIO = 1000 / 327;

export type LogoProps = {
  readonly className?: string;
  readonly height: number;
};

export function Logo({ className, height }: LogoProps) {
  return (
    <img
      className={className}
      width={height * RATIO}
      height={height}
      src={pathToLogo}
      alt="ICTSC: ICT Trouble Shooting Contest"
    />
  );
}
