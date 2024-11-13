const commentSection = document.querySelector(".comments-section");
window.addEventListener("beforeunload", () => {
  sessionStorage.setItem("scrollPositionComment", commentSection.scrollTop);
});

window.addEventListener("load", () => {
  const scrollPosition = sessionStorage.getItem("scrollPositionComment");
  if (scrollPosition) {
    commentSection.scrollTop = scrollPosition;
    sessionStorage.removeItem("scrollPositionComment"); // Supprime après usage
  }
});

// Ajuste la hauteur du textarea à l'affichage pour afficher tout le contenu sans scroll
document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll(".comment-text").forEach((textarea) => {
    textarea.style.height = "auto"; // Reset the height to auto
    textarea.style.height = textarea.scrollHeight + "px"; // Set it to the scrollHeight
  });
});

// Remove display none from edit-comment button when textarea content changes
document.querySelectorAll(".comment-text").forEach((textarea) => {
  textarea.addEventListener("input", function () {
    const editButton = this.closest(".comment").querySelector(".edit-comment");
    if (editButton) {
      editButton.style.display = "inline-block";
    }
  });
});
document.querySelectorAll("form").forEach((form) => {
  form.addEventListener("submit", (event) => {
    const textarea = form.querySelector("textarea");
    if (textarea && textarea.value.trim() === "") {
      event.preventDefault();
      alert("Comment can't be empty");
    }
  });
});