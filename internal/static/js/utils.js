(
    function () {

        function showToast(message, type = "info", width = "300px", duration = 2000, gravity = "top", position = "left") {
            const typeIcons = {
                success: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="white" width="18" height="18"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>`,
                error: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="white" width="18" height="18"><circle cx="12" cy="12" r="10" stroke="white" stroke-width="2" fill="#f44336"/><line x1="12" y1="8" x2="12" y2="13" stroke="white" stroke-width="2" stroke-linecap="round"/><circle cx="12" cy="16" r="1" fill="white"/></svg>`,
                info: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="white" width="18" height="18"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M12 20a8 8 0 100-16 8 8 0 000 16z"/></svg>`,
                warning: `<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="white" width="18" height="18"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M12 2a10 10 0 00-10 10v.34a9.93 9.93 0 001.71 5.65L12 22l8.29-4.01A9.93 9.93 0 0022 12.34V12a10 10 0 00-10-10z"/></svg>`
            };

            const typeColors = {
                success: "#4caf50", // Green
                error: "#f44336",   // Red
                info: "#2196f3",    // Blue
                warning: "#ff9800",    // Orange
            };

            let config = {
                text: `<div style="display: flex; align-items: center; gap: 10px;">${typeIcons[type] || ""}<span>${message}</span></div>`,
                duration: duration,
                gravity: gravity,
                position: position,
                stopOnFocus: true,
                escapeMarkup: false,
                style: {
                    background: typeColors[type] || "gray",
                    color: "#fff",
                    borderRadius: "8px",
                    padding: "12px 16px",
                    fontFamily: "'Segoe UI', Tahoma, Geneva, Verdana, sans-serif",
                    fontSize: "15px",
                    fontWeight: 500,
                    boxShadow: "0px 6px 12px rgba(0, 0, 0, 0.2)",
                    maxWidth: width,
                },
            };

            if (typeof options === "object") {
                config = {
                    ...config,
                    ...options,
                    text: `<div style="display: flex; align-items: center; gap: 10px;">${typeIcons[type] || ""}<span>${message}</span></div>`,
                    style: {
                        ...config.style,
                        background: typeColors[type] || config.style.background,
                        ...(options.style || {}),
                    },
                };
            }

            Toastify(config).showToast();
        }

        function dateConvertFromRawGoDates(raw) {
            const date = new Date(raw);

            let displayDate;

            if (date.getFullYear() === 1) {
                displayDate = "N/A";
            } else {
                displayDate = date.toLocaleString('en-US', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                    hour12: true,
                });
            }
            return displayDate;
        }

        function formatUnixTimestamp(ts) {
            const date = new Date(ts * 1000);
        
            if (date.getFullYear() === 1) return "N/A";
        
            return date.toLocaleString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: true,
            });
        }

        function init() {
            window.Utils = {
                showToast: showToast,
                dateConvertFromRawGoDates: dateConvertFromRawGoDates,
                formatUnixTimestamp: formatUnixTimestamp
            }
        }
        init();
    }
)()