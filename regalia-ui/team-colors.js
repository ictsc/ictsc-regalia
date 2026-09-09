const teamColors = [
  { number: "01", name: "012 U", color: "#FFE000" },
  { number: "02", name: "021 U", color: "#FF6C2F" },
  { number: "03", name: "434 U", color: "#D6CACB" },
  { number: "04", name: "Process Blue U", color: "#0083C3" },
  { number: "05", name: "356 U", color: "#397E58" },
  { number: "06", name: "Purple U", color: "#BF53B6" },
  { number: "07", name: "433 U", color: "#5B5D62" },
  { number: "08", name: "072 U", color: "#3F43AD" },
  { number: "09", name: "354 U", color: "#00A95C" },
  { number: "10", name: "3252 U", color: "#00CDC2" },
  { number: "11", name: "Violet U", color: "#7758B3" },
  { number: "12", name: "374 U", color: "#A6E35F", teamName: "KERNEL PANIC" },
  { number: "13", name: "129 U", color: "#FAAF3F" },
  { number: "14", name: "306 U", color: "#00B4E4" },
  { number: "15", name: "192 U", color: "#EE536B" },
  { number: "予備1", name: "Rubine Red U", color: "#DF547D" },
  { number: "予備2", name: "Medium Purple U", color: "#65428A" },
];

const relativeLuminance = (hex) => {
  const channels = hex.match(/[0-9a-f]{2}/gi).map((value) => {
    const channel = Number.parseInt(value, 16) / 255;
    return channel <= 0.04045 ? channel / 12.92 : ((channel + 0.055) / 1.055) ** 2.4;
  });
  return channels[0] * 0.2126 + channels[1] * 0.7152 + channels[2] * 0.0722;
};

const contrastRatio = (first, second) => {
  const light = Math.max(first, second);
  const dark = Math.min(first, second);
  return (light + 0.05) / (dark + 0.05);
};

document.querySelectorAll(".team-identity").forEach((identity) => {
  const select = document.createElement("select");

  select.className = "team-select";
  select.setAttribute("aria-label", "表示するチームカラー");

  teamColors.forEach((team) => {
    const option = document.createElement("option");
    option.value = team.number;
    option.textContent = `${team.number} / ${team.name}`;
    select.append(option);
  });

  const applyTeam = () => {
    const team = teamColors.find((item) => item.number === select.value);
    const backgroundLuminance = relativeLuminance(team.color);
    const inkLuminance = relativeLuminance("#1D252D");
    const whiteLuminance = relativeLuminance("#FFFFFF");
    const textColor = contrastRatio(backgroundLuminance, inkLuminance) >= contrastRatio(backgroundLuminance, whiteLuminance)
      ? "#1D252D"
      : "#FFFFFF";

    document.documentElement.style.setProperty("--team-color", team.color);
    document.documentElement.style.setProperty("--team-text", textColor);

    const rankingRows = document.querySelectorAll(".ranking-row[data-team]");
    const selectedRankingRow = document.querySelector(`.ranking-row[data-team="${team.number}"]`);
    rankingRows.forEach((row) => {
      const isSelected = row === selectedRankingRow;
      row.classList.toggle("is-team", isSelected);
      if (isSelected) {
        row.setAttribute("aria-current", "true");
      } else {
        row.removeAttribute("aria-current");
      }
    });

  };

  select.value = "12";
  select.addEventListener("change", applyTeam);
  identity.replaceChildren(select);
  identity.classList.add("has-selector");
  applyTeam();
});
