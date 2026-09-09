export const teamColors = [
  "#FFE000",
  "#FF6C2F",
  "#D6CACB",
  "#0083C3",
  "#397E58",
  "#BF53B6",
  "#5B5D62",
  "#3F43AD",
  "#00A95C",
  "#00CDC2",
  "#7758B3",
  "#A6E35F",
  "#FAAF3F",
  "#00B4E4",
  "#EE536B",
  "#DF547D",
  "#65428A",
] as const;
export const defaultTeamColor = "#A6E35F";
export function teamStyle(color?: string) {
  const background = teamColors.includes(color as (typeof teamColors)[number])
    ? color!
    : defaultTeamColor;
  const luminance = (hex: string) =>
    hex
      .slice(1)
      .match(/../g)!
      .map((v) => {
        const c = parseInt(v, 16) / 255;
        return c <= 0.04045 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
      })
      .reduce((sum, c, i) => sum + c * [0.2126, 0.7152, 0.0722][i]!, 0);
  const light = luminance(background),
    dark = luminance("#1D252D");
  return {
    "--team-color": background,
    "--team-text":
      (light + 0.05) / (dark + 0.05) >= 1.05 / (light + 0.05)
        ? "#1D252D"
        : "#FFFFFF",
  };
}
