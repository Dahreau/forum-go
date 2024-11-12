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
});
