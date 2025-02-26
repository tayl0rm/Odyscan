document.addEventListener("DOMContentLoaded", function () {
    const scanForm = document.getElementById("scan-form");
    const imageInput = document.getElementById("imageName");
    const scanButton = document.getElementById("scan-button");
    const progressContainer = document.getElementById("progress-container");
    const progressLog = document.getElementById("progress-log");

    scanForm.addEventListener("submit", function (event) {
        event.preventDefault();
        const imageName = imageInput.value.trim();
        if (!imageName) {
            alert("Please enter an image name.");
            return;
        }
        
        scanButton.disabled = true;
        progressContainer.style.display = "block";
        progressLog.innerHTML = "";

        // Start scanning process
        fetch("/scan", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            body: `imageName=${encodeURIComponent(imageName)}`
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Failed to start scan");
            }
            return response.text();
        })
        .then(() => {
            listenForProgress();
        })
        .catch(error => {
            progressLog.innerHTML += `<p class='error'>Error: ${error.message}</p>`;
            scanButton.disabled = false;
        });
    });

    function listenForProgress() {
        const eventSource = new EventSource("/scan/progress");

        eventSource.onmessage = function (event) {
            const logEntry = document.createElement("p");
            logEntry.textContent = event.data;
            progressLog.appendChild(logEntry);
            progressLog.scrollTop = progressLog.scrollHeight;

            if (event.data.includes("✅")) {
                eventSource.close();
                scanButton.disabled = false;
            }
        };

        eventSource.onerror = function () {
            eventSource.close();
            progressLog.innerHTML += `<p class='error'>Error: Lost connection to server.</p>`;
            scanButton.disabled = false;
        };
    }
});
