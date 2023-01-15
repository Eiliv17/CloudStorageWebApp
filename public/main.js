let input = document.getElementById("file");
const filename = document.querySelector("span.filename");

input.addEventListener("change", ()=>{
    let fileinput = document.querySelector("input[type=file]").files[0];

    filename.textContent = fileinput.name;
})