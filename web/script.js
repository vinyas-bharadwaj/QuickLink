document.addEventListener('DOMContentLoaded', () => {
  const form = document.getElementById('shorten-form');
  const urlInput = document.getElementById('url-input');
  const submitBtn = document.getElementById('submit-btn');
  const errorAlert = document.getElementById('error-alert');
  const resultSection = document.getElementById('result-section');
  const resultOriginal = document.getElementById('result-original');
  const resultShortLink = document.getElementById('result-short-link');
  const copyBtn = document.getElementById('copy-btn');
  const visitBtn = document.getElementById('visit-btn');
  const historyList = document.getElementById('history-list');
  const historyEmpty = document.getElementById('history-empty');
  const historyCount = document.getElementById('history-count');
  const clearHistoryBtn = document.getElementById('clear-history-btn');
  const toast = document.getElementById('toast');

  const STORAGE_KEY = 'quicklink_recent_urls';

  renderHistory();

  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();

    let rawUrl = urlInput.value.trim();
    if (!rawUrl) return showError('Please enter a URL to shorten.');
    if (!/^https?:\/\//i.test(rawUrl)) rawUrl = 'https://' + rawUrl;

    submitBtn.disabled = true;
    submitBtn.textContent = 'Shortening...';

    try {
      const res = await fetch('/shorten', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: new URLSearchParams({ url: rawUrl })
      });

      if (!res.ok) throw new Error(await res.text() || 'Failed to shorten URL');
      const data = await res.json();
      
      const fullShortUrl = `${window.location.origin}/${data.short_url}`;
      
      resultOriginal.textContent = data.original_url;
      resultShortLink.textContent = fullShortUrl;
      resultShortLink.href = fullShortUrl;
      visitBtn.href = fullShortUrl;
      resultSection.classList.remove('hidden');

      saveHistory({ id: Date.now(), shortUrl: fullShortUrl, originalUrl: data.original_url });
      urlInput.value = '';
    } catch (err) {
      showError(err.message);
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Shorten URL';
    }
  });

  copyBtn.addEventListener('click', () => copyToClipboard(resultShortLink.href));
  clearHistoryBtn.addEventListener('click', () => {
    localStorage.removeItem(STORAGE_KEY);
    renderHistory();
    showToast('History cleared');
  });

  function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => showToast('Copied to clipboard!'));
  }

  function showToast(msg) {
    toast.textContent = msg;
    toast.classList.remove('hidden');
    setTimeout(() => toast.classList.add('hidden'), 2200);
  }

  function showError(msg) {
    errorAlert.textContent = msg;
    errorAlert.classList.remove('hidden');
  }

  function hideError() {
    errorAlert.classList.add('hidden');
  }

  function getHistory() {
    return JSON.parse(localStorage.getItem(STORAGE_KEY) || '[]');
  }

  function saveHistory(item) {
    const history = [item, ...getHistory().filter(h => h.shortUrl !== item.shortUrl)].slice(0, 15);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history));
    renderHistory();
  }

  function renderHistory() {
    const history = getHistory();
    historyCount.textContent = history.length;

    if (!history.length) {
      historyEmpty.classList.remove('hidden');
      historyList.classList.add('hidden');
      clearHistoryBtn.classList.add('hidden');
      return;
    }

    historyEmpty.classList.add('hidden');
    historyList.classList.remove('hidden');
    clearHistoryBtn.classList.remove('hidden');

    historyList.innerHTML = history.map(item => `
      <div class="history-item">
        <div class="history-urls">
          <a class="history-short" href="${item.shortUrl}" target="_blank">${item.shortUrl}</a>
          <span class="history-original" title="${item.originalUrl}">${item.originalUrl}</span>
        </div>
        <div class="history-actions">
          <button class="btn btn-secondary" onclick="navigator.clipboard.writeText('${item.shortUrl}').then(() => { const t = document.getElementById('toast'); t.textContent='Copied to clipboard!'; t.classList.remove('hidden'); setTimeout(()=>t.classList.add('hidden'), 2200); })">Copy</button>
          <button class="btn btn-ghost" onclick="removeHistory(${item.id})">✕</button>
        </div>
      </div>
    `).join('');
  }

  window.removeHistory = (id) => {
    const history = getHistory().filter(h => h.id !== id);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history));
    renderHistory();
  };
});
