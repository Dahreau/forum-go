document.addEventListener("DOMContentLoaded", () => {
  function setActiveTab() {
    const tabs = document.querySelectorAll(".tab");
    const path = window.location.pathname;

    // Définir l'onglet actif en fonction de l'URL
    tabs.forEach((tab) => tab.classList.remove("active")); // Réinitialiser l'onglet actif

    if (path === "/created") {
      tabs[1].classList.add("active");
    } else if (path === "/liked") {
      tabs[2].classList.add("active");
    } else {
      tabs[0].classList.add("active"); // Par défaut, "All posts"
    }
  }

  setActiveTab();
  sortOptions();
});

const select = document.getElementById("filter-category");
const selectedCategories = document.getElementById("selected-categories");
const originalOptions = Array.from(select.options);

function addCategory() {
  const selectedOption = select.options[select.selectedIndex];
  if (selectedOption) {
    const categoryElement = document.createElement("div");
    categoryElement.className = "category-box";
    categoryElement.innerHTML = `<span class="remove-btn">×</span> ${selectedOption.text}`;
    categoryElement.id = `selected-${selectedOption.value}`;
    categoryElement.onclick = function () {
      removeCategory(selectedOption.value, selectedOption.text);
    };
    selectedCategories.appendChild(categoryElement);
    select.remove(select.selectedIndex);
  }
}

function removeCategory(value, text) {
  const categoryElement = document.getElementById(`selected-${value}`);
  if (categoryElement) {
    categoryElement.remove();
  }
  const option = document.createElement("option");
  option.value = value;
  option.text = text;
  select.add(option);
  sortOptions();
}
select.onchange = addCategory;

function sortOptions() {
  const options = Array.from(select.options).slice(1); // Ignore first option
  options.sort((a, b) => a.text.localeCompare(b.text));
  select.innerHTML = '<option value="">Select one or more categories</option>';
  options.forEach((option) => select.add(option));
}

const btnResetFilters = document.getElementById("btn-reset-filters");
if (btnResetFilters) {
  btnResetFilters.onclick = function () {
    originalOptions.forEach((option) => select.add(option));
    selectedCategories.innerHTML = "";
    select.selectedIndex = 0;
    sortOptions();
  };
}
