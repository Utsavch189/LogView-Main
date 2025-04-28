let refreshInterval = 10;
let refreshTimer;
let isAutoRefreshEnabled = false;
const refreshIndicator = document.getElementById('refreshIndicator');

function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function updateTimer(seconds) {
    if(isAutoRefreshEnabled){
        refreshIndicator.classList.add('bg-green-500');
        document.getElementById('timer').textContent = seconds + 's';
    }
    else{
        refreshIndicator.classList.remove('bg-green-500');
        document.getElementById('timer').textContent = '';
    }
}

function toggleAutoRefresh() {
    isAutoRefreshEnabled = !isAutoRefreshEnabled;
    if (isAutoRefreshEnabled) {
        refreshIndicator.classList.add('bg-green-500');
        refreshCycle();
    } else {
        refreshIndicator.classList.remove('bg-green-500')
        clearInterval(refreshTimer);
        updateTimer('-');
    }
}

function updateRefreshInterval(value) {
    refreshInterval = parseInt(value) || 5;
    if (isAutoRefreshEnabled) {
        clearInterval(refreshTimer);
        refreshCycle();
    }
}

function refreshCycle() {
    let countdown = refreshInterval;

    if (isAutoRefreshEnabled===true){
        function tick() {
            updateTimer(countdown);
            countdown--;

            if (countdown < 0) {
                countdown = refreshInterval;
                getLogs();
            }
        }

        getLogs();
        tick();
        refreshTimer = setInterval(tick, 1000);
    }
}

refreshCycle();