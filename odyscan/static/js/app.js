document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("scanForm");
    const resultDiv = document.getElementById("result");

    form.addEventListener("submit", function (event) {
        event.preventDefault();
        
        const imageName = document.getElementById("imageName").value;
        if (!imageName) {
            resultDiv.innerHTML = `<p style="color: red;">Please enter an image name.</p>`;
            return;
        }

        fetch("/scan", {
            method: "POST",
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            body: `imageName=${encodeURIComponent(imageName)}`
        })
        .then(response => response.text())
        .then(data => {
            resultDiv.innerHTML = `<p>${data}</p>`;
        })
        .catch(error => {
            resultDiv.innerHTML = `<p style="color: red;">Error: ${error.message}</p>`;
        });
    });
});
