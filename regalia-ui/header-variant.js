const pageParams = new URLSearchParams(window.location.search);
const problemVariant = pageParams.get("problem");
const legacySimpleVariant = pageParams.get("header") === "simple";

document.body.classList.add("simple-header-version");

if (problemVariant === "a04" || legacySimpleVariant) {
  document.title = "A04 無線LANに接続できない / ICTSC 2026 Regalia";

  const currentProblem = document.querySelector(".problem-rail a.is-active");
  currentProblem?.classList.remove("is-active");
  currentProblem?.removeAttribute("aria-current");
  currentProblem?.setAttribute("aria-label", currentProblem.getAttribute("aria-label").replace(" 現在の問題", ""));

  const a04Problem = [...document.querySelectorAll(".problem-rail a")].find((link) =>
    link.querySelector("b")?.textContent.trim() === "A04"
  );
  a04Problem?.classList.add("is-active");
  a04Problem?.setAttribute("aria-current", "page");
  a04Problem?.setAttribute("aria-label", "A04 採点中 現在の問題 再提出まで18分");

  document.querySelector(".hero-copy h1").textContent = "A04:無線LANに接続できない.";
  document.querySelector(".problem-score-summary strong").textContent = "—";
  document.querySelector(".problem-score-summary span").textContent = "/ 200点満点";
  const problemAvailability = document.querySelector(".problem-cooldown");
  problemAvailability.classList.remove("is-answer-closed");
  const availabilityLabel = document.createElement("span");
  const availabilityValue = document.createElement("strong");
  const availabilityUnit = document.createElement("small");
  availabilityLabel.textContent = "再提出まで";
  availabilityValue.textContent = "18";
  availabilityUnit.textContent = "分";
  problemAvailability.replaceChildren(availabilityLabel, availabilityValue, availabilityUnit);

  const problemParagraphs = document.querySelectorAll(".reading-copy p");
  problemParagraphs[0].textContent = "会場内の端末から指定された無線LANへ接続できません。アクセスポイントと端末の状態を確認し、原因を特定してください。";
  problemParagraphs[1].textContent = "必要に応じて設定を修正し、接続と疎通が回復したことを確認してください。変更内容と確認に使ったコマンドを回答欄に記入してください。";

  const environmentValues = document.querySelectorAll(".environment-list div");
  const a04Environment = [
    ["ap-a", "10.30.1.1"],
    ["client-a", "10.30.1.20"],
    ["SSID", "ICTSC-WIFI"],
    ["周波数帯", "5GHz"],
  ];
  environmentValues.forEach((row, index) => {
    row.querySelector("dt").textContent = a04Environment[index][0];
    row.querySelector("dd").textContent = a04Environment[index][1];
  });
}
