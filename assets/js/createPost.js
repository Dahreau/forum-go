let button = document.querySelector(".upload-button");
let deleteButton = document.querySelector("#deleteFile");
console.log(button);
button.addEventListener("change", (e) => {
  fileNameSpan = document.querySelector(".file-name");
  fileNameSpan.innerHTML = e.target.files[0].name;
  console.log(e.target.files);
  if (e.target.files.length > 0) {
    deleteButton.style.display = "block";
  }
});

deleteButton.addEventListener("click", (e) => {
  fileNameSpan.innerHTML = "No file chosen";
  button.value = "";
  deleteButton.style.display = "none";
});
