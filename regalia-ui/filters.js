const statusFilter = document.querySelector("#status-filter");
const dayFilter = document.querySelector("#day-filter");
const categoryFilter = document.querySelector("#category-filter");
const problemRows = document.querySelectorAll("a.problem-row");
const problemTables = document.querySelectorAll(".problem-table");

function updateProblemRows() {
  const status = statusFilter.value;
  const day = dayFilter.value;
  const category = categoryFilter.value;

  problemRows.forEach((row) => {
    const matchesStatus = status === "all" || row.dataset.status === status;
    const matchesDay = day === "all" || row.dataset.day === day;
    const matchesCategory = category === "all" || row.dataset.category === category;
    row.hidden = !(matchesStatus && matchesDay && matchesCategory);
  });

  problemTables.forEach((table) => {
    const rows = table.querySelectorAll("a.problem-row");
    table.hidden = !Array.from(rows).some((row) => !row.hidden);
  });
}

statusFilter.addEventListener("change", updateProblemRows);
dayFilter.addEventListener("change", updateProblemRows);
categoryFilter.addEventListener("change", updateProblemRows);
