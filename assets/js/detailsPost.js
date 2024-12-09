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
  document.querySelectorAll("textarea").forEach((textarea) => {
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
// Remove display none from edit-post button when textarea content changes
let postText = document.querySelector(".post-text");

postText.addEventListener("input", function () {
  const editButton = this.closest(".post-content").querySelector(".edit-post");
  if (editButton) {
    editButton.style.display = "inline-block";
  }
});
document.querySelectorAll("form").forEach((form) => {
  form.addEventListener("submit", (event) => {
    const textarea = form.querySelector("textarea");
    if (textarea && textarea.value.trim() === "") {
      event.preventDefault();
      alert("Field can't be empty");
    }
  });
});
