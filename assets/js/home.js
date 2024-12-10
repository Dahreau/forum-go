const select = document.getElementById("filter-category");
const selectedCategories = document.getElementById("selected-categories");
const originalOptions = Array.from(select.options);

document.addEventListener("DOMContentLoaded", () => {
  function setActiveTab() {
    const tabs = document.querySelectorAll(".tab");
    const path = window.location.pathname;

    if (tabs.length > 0) {
      // Set the active tab based on the URL
      tabs.forEach((tab) => tab.classList.remove("active")); // Reset the active tab

      if (path === "/created") {
        tabs[1].classList.add("active");
      } else if (path === "/liked") {
        tabs[2].classList.add("active");
      } else {
        tabs[0].classList.add("active"); // Default to "All posts"
      }
    }
  }
  sortOptions();
  setActiveTab();
});

function addCategory() {
  // Add the selected category to the selected categories
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
    sortOptions();
  }
}

function removeCategory(value, text) {
  // Remove the category from the selected categories and add it back to the select element.
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
  // Sort the options in the select element
  const options = Array.from(select.options).slice(1); // Ignore first option
  options.sort((a, b) => a.text.localeCompare(b.text));
  select.innerHTML = '<option value="">Select one or more categories</option>';
  options.forEach((option) => select.add(option));
  displayPosts();
}

const btnResetFilters = document.getElementById("btn-reset-filters");
// Reset all selected categories
if (btnResetFilters) {
  btnResetFilters.onclick = function () {
    Array.from(selectedCategories.children).forEach((category) => {
      removeCategory(
        category.id.replace("selected-", ""),
        category.textContent.replace("×", "").trim()
      );
    });
  };
}

function displayPosts() {
  // Get all posts and categories to filter by  (post-item and category-box)
  const posts = document.querySelectorAll(".post-item");
  const arrayCategories = Array.from(selectedCategories.children).map(
    (category) => category.textContent.replace("×", "").trim()
  );

  posts.forEach((post) => {
    const categories = Array.from(post.querySelectorAll(".category-box")).map(
      (category) => category.textContent
    );
    const isDisplayed = arrayCategories.every((category) =>
      categories.includes(category)
    );
    post.style.display = isDisplayed ? "block" : "none";
  });
}
