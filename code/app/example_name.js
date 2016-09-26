chrome.runtime.getBackgroundPage(function(bg) {
    // Relevant function at the background page. In your specific example:
    bg.clearCache();
});
