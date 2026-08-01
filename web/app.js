const form = document.getElementById("shorten-form");
const urlInput = document.getElementById("url-input");
const result = document.getElementById("result");
const error = document.getElementById("error");

form.addEventListener("submit", async (event) => {
  event.preventDefault();
  result.style.display = "none";
  result.textContent = "";
  error.textContent = "";

  const formData = new URLSearchParams();
  formData.set("url", urlInput.value.trim());

  try {
    const response = await fetch("/shorten", {
      method: "POST",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: formData.toString(),
    });

    if (!response.ok) {
      throw new Error(await response.text());
    }

    const data = await response.json();
    const shortLink = `${window.location.origin}/${data.short_url}`;

    result.innerHTML = `Short URL: <a href="${shortLink}" target="_blank" rel="noopener noreferrer">${shortLink}</a>`;
    result.style.display = "block";
  } catch (e) {
    error.textContent = e.message || "Something went wrong while shortening the URL.";
  }
});
