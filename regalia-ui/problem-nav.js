const problemRailToggle = document.querySelector(".problem-rail-toggle");
const problemRailToggleLabel = problemRailToggle?.querySelector(".visually-hidden");
const problemRail = document.querySelector("#problem-rail");
const detailLayout = document.querySelector(".detail-layout");

problemRailToggle?.addEventListener("click", () => {
  const willHide = !problemRail.hidden;
  problemRail.hidden = willHide;
  detailLayout.classList.toggle("is-rail-hidden", willHide);
  problemRailToggle.setAttribute("aria-expanded", String(!willHide));
  const label = willHide ? "問題リンクを表示" : "問題リンクを隠す";
  problemRailToggle.setAttribute("aria-label", label);
  problemRailToggle.setAttribute("title", label);
  problemRailToggleLabel.textContent = label;
});
