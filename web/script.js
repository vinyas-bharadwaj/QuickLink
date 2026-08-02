/**
 * QuickLink Frontend Core Script
 * Fully native Vanilla JS implementation for high performance & reliability.
 */

document.addEventListener('DOMContentLoaded', () => {
  // DOM Element References
  const form = document.getElementById('shorten-form');
  const urlInput = document.getElementById('url-input');
  const submitBtn = document.getElementById('submit-btn');
  const btnLabel = document.getElementById('btn-label');
  const btnSpinner = document.getElementById('btn-spinner');
  const errorAlert = document.getElementById('error-alert');
  const errorText = document.getElementById('error-text');
  
  // Result Section
  const resultSection = document.getElementById('result-section');
  const resultOriginal = document.getElementById('result-original');
  const resultShortLink = document.getElementById('result-short-link');
  const copyBtn = document.getElementById('copy-btn');
  const copyIcon = document.getElementById('copy-icon');
  const checkIcon = document.getElementById('check-icon');
  const copyText = document.getElementById('copy-text');
  const visitBtn = document.getElementById('visit-btn');
  
  // History Section
  const historyList = document.getElementById('history-list');
  const historyEmpty = document.getElementById('history-empty');
  const historyCount = document.getElementById('history-count');
  const clearHistoryBtn = document.getElementById('clear-history-btn');
  
  // Toast Element
  const toast = document.getElementById('toast');
  const toastText = document.getElementById('toast-text');

  // Constants
  const STORAGE_KEY = 'quicklink_recent_urls';
  const MAX_HISTORY = 20;

  // Initial setup: load history from storage
  loadAndRenderHistory();

  // Focus input on page load for immediate utility
  if (urlInput) {
    urlInput.focus();
  }

  // Handle Form Submission
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();

    let rawUrl = urlInput.value.trim();
    if (!rawUrl) {
      showError('Please paste or type a URL to shorten.');
      urlInput.focus();
      return;
    }

    // Auto-prepend HTTPS protocol if omitted by user
    if (!/^https?:\/\//i.test(rawUrl) && !/^ftp:\/\//i.test(rawUrl)) {
      rawUrl = 'https://' + rawUrl;
      urlInput.value = rawUrl;
    }

    // Validate URL syntax
    if (!isValidUrl(rawUrl)) {
      showError('The provided address does not look like a valid web link.');
      urlInput.focus();
      return;
    }

    // Initiate API request
    setLoadingState(true);
    
    try {
      const response = await fetch('/shorten', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: new URLSearchParams({ url: rawUrl })
      });

      if (!response.ok) {
        let errorMessage = 'Failed to shorten URL. Please try again.';
        try {
          const errData = await response.text();
          if (errData) errorMessage = errData;
        } catch (err) {}
        throw new Error(errorMessage);
      }

      const data = await response.json();
      
      // Parse payload: data contains message, short_url (the code), and original_url
      const shortCode = data.short_url;
      const fullShortUrl = `${window.location.origin}/${shortCode}`;
      
      // Display result cleanly
      displayResult(data.original_url, fullShortUrl);
      
      // Save item in localStorage history
      saveToHistory({
        id: Date.now(),
        code: shortCode,
        shortUrl: fullShortUrl,
        originalUrl: data.original_url,
        created: new Date().toISOString()
      });

      // Reset input
      urlInput.value = '';
      
    } catch (err) {
      console.error('Shorten Error:', err);
      showError(err.message || 'An unexpected server error occurred.');
    } finally {
      setLoadingState(false);
    }
  });

  // Handle Result Copy Button
  copyBtn.addEventListener('click', () => {
    const urlToCopy = resultShortLink.href;
    copyToClipboard(urlToCopy, () => {
      // Button Feedback Animation
      copyIcon.classList.add('hidden');
      copyIcon.setAttribute('hidden', '');
      checkIcon.classList.remove('hidden');
      checkIcon.removeAttribute('hidden');
      copyText.textContent = 'Copied!';
      copyBtn.classList.add('border-active');

      setTimeout(() => {
        copyIcon.classList.remove('hidden');
        copyIcon.removeAttribute('hidden');
        checkIcon.classList.add('hidden');
        checkIcon.setAttribute('hidden', '');
        copyText.textContent = 'Copy';
        copyBtn.classList.remove('border-active');
      }, 2000);
    });
  });

  // Handle Clear History
  clearHistoryBtn.addEventListener('click', () => {
    if (confirm('Clear all saved URL shortcuts from this browser?')) {
      localStorage.removeItem(STORAGE_KEY);
      loadAndRenderHistory();
      showToast('History cleared.');
    }
  });

  /* ==========================================================================
     Helper & UI Functions
     ========================================================================== */

  function displayResult(original, shortUrl) {
    resultOriginal.textContent = original;
    resultOriginal.title = original;
    
    resultShortLink.textContent = shortUrl;
    resultShortLink.href = shortUrl;
    visitBtn.href = shortUrl;
    
    // Reveal result card with animation
    resultSection.classList.remove('hidden');
    resultSection.removeAttribute('hidden');
    resultSection.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  }

  function setLoadingState(isLoading) {
    if (isLoading) {
      submitBtn.disabled = true;
      urlInput.disabled = true;
      btnLabel.classList.add('hidden');
      btnLabel.setAttribute('hidden', '');
      btnSpinner.classList.remove('hidden');
      btnSpinner.removeAttribute('hidden');
    } else {
      submitBtn.disabled = false;
      urlInput.disabled = false;
      btnLabel.classList.remove('hidden');
      btnLabel.removeAttribute('hidden');
      btnSpinner.classList.add('hidden');
      btnSpinner.setAttribute('hidden', '');
      urlInput.focus();
    }
  }

  function showError(msg) {
    errorText.textContent = msg;
    errorAlert.classList.remove('hidden');
    errorAlert.removeAttribute('hidden');
  }

  function hideError() {
    errorAlert.classList.add('hidden');
    errorAlert.setAttribute('hidden', '');
  }

  function isValidUrl(string) {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  }

  function copyToClipboard(text, onSuccess) {
    if (navigator.clipboard && window.isSecureContext) {
      navigator.clipboard.writeText(text).then(() => {
        if (onSuccess) onSuccess();
        showToast('Link copied to clipboard!');
      }).catch(err => {
        console.error('Clipboard API write failed:', err);
        fallbackCopy(text, onSuccess);
      });
    } else {
      fallbackCopy(text, onSuccess);
    }
  }

  function fallbackCopy(text, onSuccess) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.position = 'fixed';
    textArea.style.left = '-999999px';
    textArea.style.top = '-999999px';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    try {
      document.execCommand('copy');
      if (onSuccess) onSuccess();
      showToast('Link copied to clipboard!');
    } catch (err) {
      showError('Failed to copy to clipboard. Please copy manually.');
    }
    textArea.remove();
  }

  function showToast(msg, duration = 2500) {
    toastText.textContent = msg;
    toast.classList.remove('hidden');
    toast.removeAttribute('hidden');
    
    clearTimeout(toast._timer);
    toast._timer = setTimeout(() => {
      toast.classList.add('hidden');
      toast.setAttribute('hidden', '');
    }, duration);
  }

  /* ==========================================================================
     Local Storage History Management
     ========================================================================== */

  function getHistory() {
    try {
      const data = localStorage.getItem(STORAGE_KEY);
      return data ? JSON.parse(data) : [];
    } catch (err) {
      console.error('Error reading localStorage:', err);
      return [];
    }
  }

  function saveToHistory(newItem) {
    let history = getHistory();
    
    // Check if original URL or code already exists and remove previous entry to bump to top
    history = history.filter(item => item.code !== newItem.code);
    
    // Add new item to beginning
    history.unshift(newItem);
    
    // Enforce max count
    if (history.length > MAX_HISTORY) {
      history = history.slice(0, MAX_HISTORY);
    }

    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(history));
    } catch (err) {
      console.error('Error writing to localStorage:', err);
    }

    loadAndRenderHistory();
  }

  function removeHistoryItem(id) {
    let history = getHistory();
    history = history.filter(item => item.id !== id);
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(history));
    } catch (err) {
      console.error('Error updating localStorage:', err);
    }
    loadAndRenderHistory();
    showToast('Item removed from history.');
  }

  function loadAndRenderHistory() {
    const history = getHistory();
    historyCount.textContent = history.length.toString();

    if (history.length === 0) {
      historyEmpty.classList.remove('hidden');
      historyEmpty.removeAttribute('hidden');
      historyList.classList.add('hidden');
      historyList.setAttribute('hidden', '');
      clearHistoryBtn.classList.add('hidden');
      clearHistoryBtn.setAttribute('hidden', '');
      return;
    }

    historyEmpty.classList.add('hidden');
    historyEmpty.setAttribute('hidden', '');
    historyList.classList.remove('hidden');
    historyList.removeAttribute('hidden');
    clearHistoryBtn.classList.remove('hidden');
    clearHistoryBtn.removeAttribute('hidden');
    
    // Render list securely using DocumentFragment and safe property assignments
    historyList.innerHTML = '';
    const fragment = document.createDocumentFragment();

    history.forEach(item => {
      const itemEl = document.createElement('div');
      itemEl.className = 'history-item';
      
      const contentDiv = document.createElement('div');
      contentDiv.className = 'history-item-content';
      
      const shortLink = document.createElement('a');
      shortLink.className = 'history-short-link';
      shortLink.href = item.shortUrl;
      shortLink.target = '_blank';
      shortLink.rel = 'noopener noreferrer';
      shortLink.textContent = item.shortUrl;
      
      const originalLink = document.createElement('a');
      originalLink.className = 'history-original-link';
      originalLink.href = item.originalUrl;
      originalLink.target = '_blank';
      originalLink.rel = 'noopener noreferrer';
      originalLink.textContent = item.originalUrl;
      originalLink.title = item.originalUrl;
      
      contentDiv.appendChild(shortLink);
      contentDiv.appendChild(originalLink);
      
      const actionsDiv = document.createElement('div');
      actionsDiv.className = 'history-item-actions';
      
      // Copy Button
      const rowCopyBtn = document.createElement('button');
      rowCopyBtn.className = 'btn btn-icon';
      rowCopyBtn.title = 'Copy shortened URL';
      rowCopyBtn.setAttribute('aria-label', 'Copy shortened URL');
      rowCopyBtn.innerHTML = `
        <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
        </svg>
      `;
      rowCopyBtn.addEventListener('click', () => {
        copyToClipboard(item.shortUrl);
      });
      
      // Delete Button
      const rowDeleteBtn = document.createElement('button');
      rowDeleteBtn.className = 'btn btn-icon';
      rowDeleteBtn.title = 'Remove from history';
      rowDeleteBtn.setAttribute('aria-label', 'Remove from history');
      rowDeleteBtn.innerHTML = `
        <svg viewBox="0 0 24 24" width="15" height="15" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="3 6 5 6 21 6"></polyline>
          <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
          <line x1="10" y1="11" x2="10" y2="17"></line>
          <line x1="14" y1="11" x2="14" y2="17"></line>
        </svg>
      `;
      rowDeleteBtn.addEventListener('click', () => {
        removeHistoryItem(item.id);
      });

      actionsDiv.appendChild(rowCopyBtn);
      actionsDiv.appendChild(rowDeleteBtn);
      
      itemEl.appendChild(contentDiv);
      itemEl.appendChild(actionsDiv);
      
      fragment.appendChild(itemEl);
    });

    historyList.appendChild(fragment);
  }
});
