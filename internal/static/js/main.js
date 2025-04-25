let checkcked_log_levels = [
    { "level": "info", "id": "info-check", "checked": true },
    { "level": "warning", "id": "warning-check", "checked": true },
    { "level": "error", "id": "error-check", "checked": true },
    { "level": "debug", "id": "debug-check", "checked": true },
];

let dates = {
    from: null,
    to: null
};

let currentPage = 1;
const itemsPerPage = 20;

function openLogModal(log) {
    document.getElementById("logDetails").innerText = JSON.stringify(log, null, 2);
    document.getElementById("logModal").classList.remove("hidden");
}

function closeLogModal() {
    document.getElementById("logModal").classList.add("hidden");
}

function openProjectModal() {
    document.getElementById("projectModal").classList.remove("hidden");
}

function closeProjectModal() {
    document.getElementById("projectModal").classList.add("hidden");
}

async function getLogs(project_name = "") {
    try {
        if (project_name === "") {
            const urlParams = new URLSearchParams(window.location.search);
            project_name = urlParams.get('project');

            if (!project_name) {
                return;
            }
        }
        else {
            const url = new URL(window.location.href);
            url.searchParams.set('project', project_name);
            window.history.pushState({}, '', url);
        }

        const res = await fetch(`/api/logs/${project_name}/apply-filters/get-all?page=${currentPage}&page_size=${itemsPerPage}`, {
            method: "POST",
            body: JSON.stringify({
                "loglevels": checkcked_log_levels,
                "dates": dates
            })
        })
        const data = await res.json();
        renderLogs(data?.logs, data?.count, data?.info_count, data?.warn_count, data?.error_count, data?.debug_count,data?.paginate_count);
    } catch (error) {
        console.log(error)
        window.Utils.showToast("Something is wrong!", "error")
    }
}

async function createProject() {
    try {
        const name = document.getElementById("projectName").value.trim();
        if (name) {
            const res = await fetch("/api/project/create", {
                method: "POST",
                body: JSON.stringify({
                    "project_name": name
                })
            })
            const data = await res.json();
            window.Utils.showToast(`Project ${data?.project_name} is created!`, "success");
            await getProjects();
            closeLogModal();
        }
        else {
            window.Utils.showToast("Enter a Project Name!", "error");
        }
    } catch (error) {
        console.log(error)
        window.Utils.showToast("Something is wrong!", "error")
        closeProjectModal();
    }
    finally {
        document.getElementById("projectName").value = "";
    }
}

async function getProjects() {
    try {
        const res = await fetch("/api/project/get-all")
        const data = await res.json();
        renderProjects(data)
    } catch (error) {
        window.Utils.showToast("Something is wrong!", "error")
    }
}

function renderLogs(logs, total_logs, info_count, warn_count, error_count, debug_count, paginate_count) {
    const container = document.getElementById("logContainer");
    container.innerHTML = '';

    document.getElementById("total_logs_count").innerText = total_logs;
    document.getElementById("info_logs_count").innerText = info_count;
    document.getElementById("error_logs_count").innerText = error_count;
    document.getElementById("warn_logs_count").innerText = warn_count;
    document.getElementById("debug_logs_count").innerText = debug_count;

    if (!logs || !logs?.length || total_logs === 0) {
        return;
    }

    const totalPages = Math.ceil(paginate_count / itemsPerPage);
    currentPage = Math.min(currentPage, totalPages);

    container.innerHTML = logs.map(log => `
        <tr class="border-b hover:bg-gray-50">
          <td class="px-4 py-2">${window.Utils.dateConvertFromRawGoDates(log.created_at)}</td>
          <td class="px-4 py-2 text-${log.level === "ERROR" ? "red" : log.level === "DEBUG" ? "gray" : log.level === "WARNING" ? "yellow" : "green"}-500 font-semibold">${log.level}</td>
          <td class="px-4 py-2">${log.logger}</td>
          <td class="px-4 py-2">${log.message}</td>
          <td class="px-4 py-2">
            <button onclick='openLogModal(${JSON.stringify(log)})' class="text-blue-600 hover:underline">View</button>
          </td>
        </tr>
      `).join("");

    document.getElementById("pageInfo").textContent = currentPage;
    document.getElementById("totalPages").textContent = totalPages;
    document.getElementById("totalCount").textContent = paginate_count;

    document.getElementById("prevBtn").disabled = currentPage <= 1;
    document.getElementById("nextBtn").disabled = currentPage >= totalPages;

}

function nextPage() {
    currentPage++;
    getLogs();
}

function prevPage() {
    if (currentPage > 1) currentPage--;
    getLogs();
}

function toggleCollapsible(element) {
    const content = element.nextElementSibling;
    const arrow = element.querySelector('span:last-child');

    content.classList.toggle('hidden');
    arrow.style.transform = content.classList.contains('hidden') ? '' : 'rotate(180deg)';
}

function renderProjects(projects) {

    const list = document.getElementById("projectList");
    list.innerHTML = "";

    if (!projects || !projects?.length) {
        window.Utils.showToast("No Projects Found. Please Create one!", "warning");
        return;
    }
    projects.map((v) => {
        const elm = `
            <div class="flex items-center justify-center p-4">
                <div class="w-full max-w-md">
                  <div id="collapsible" class="bg-white rounded-lg shadow-md">
                    <div 
                      class="cursor-pointer px-4 py-3 flex items-center justify-between bg-white hover:bg-gray-50 transition-colors"
                      onclick="toggleCollapsible(this)"
                    >
                      <span class="font-medium text-gray-800">${v?.project_name}</span>
                      <span class="transform transition-transform">▼</span>
                    </div>

                    <div class="hidden px-4 py-3 border-t border-gray-200">
                      <div class="space-y-3">
                        <div>
                          <label class="block text-sm font-medium text-gray-600 mb-1">
                            Source Token
                          </label>
                          <div class="text-sm text-gray-800 font-mono bg-gray-50 p-2 rounded">
                            ${v?.source_token}
                          </div>
                        </div>
                        <button 
                          class="w-full bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition-colors"
                          onclick="getLogs('${v?.project_name}')"
                        >
                          View Log
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
            </div>
        `
        list.innerHTML += elm;
    })
}

function set_checking_loglevels() {
    checkcked_log_levels.forEach(v => {
        if (v.checked === true) {
            document.getElementById(v.id).setAttribute('checked', true);
        }
        else {
            document.getElementById(v.id).removeAttribute('checked');
        }
    })
}

document.querySelectorAll(".loglevel_checks").forEach(i => {
    i.addEventListener("change", (e) => {
        const checked = e.target.getAttribute("checked");

        if (checked === 'true') {
            console.log(e.target.getAttribute("id"))

            checkcked_log_levels.forEach(v => {
                if (v.id === e.target.getAttribute("id")) {
                    v.checked = false;
                    document.getElementById(v.id).removeAttribute('checked');
                }
            })
        }
        else {
            checkcked_log_levels.forEach(v => {
                if (v.id === e.target.getAttribute("id")) {
                    v.checked = true;
                    document.getElementById(v.id).setAttribute('checked', true);
                }
            })
        }
        getLogs();
    })
})

function validateDates() {
    const fromDate = document.getElementById('fromDate');
    const toDate = document.getElementById('toDate');
    const selectedRange = document.getElementById('selectedRange');

    if (fromDate.value && toDate.value) {
        const from = new Date(fromDate.value);
        const to = new Date(toDate.value);

        if (from > to) {
            toDate.value = fromDate.value;
        }

        // Format dates for display
        const formatDate = (date) => {
            return new Intl.DateTimeFormat('en-US', {
                weekday: 'short',
                year: 'numeric',
                month: 'short',
                day: 'numeric'
            }).format(date);
        };

        selectedRange.textContent = `${formatDate(from)} - ${formatDate(to)}`;
    } else {
        selectedRange.textContent = 'No dates selected';
    }
}

document.getElementById("dateSearch").addEventListener("click", async (e) => {
    const fromDate = document.getElementById('fromDate').value;
    const toDate = document.getElementById('toDate').value;

    if (!fromDate || !toDate) {
        window.Utils.showToast("Provide Date Ranges!", "error");
        return;
    }
    dates['from'] = new Date(fromDate).toISOString();
    dates['to'] = new Date(toDate).toISOString();
    await getLogs();
})

document.getElementById("resetDateSearch").addEventListener("click", async (e) => {
    document.getElementById('fromDate').value = "";
    document.getElementById('toDate').value = "";

    dates['from'] = null;
    dates['to'] = null;
    await getLogs();
})


document.getElementById("log-download-btn").addEventListener("click", async (e) => {
    const urlParams = new URLSearchParams(window.location.search);
    project_name = urlParams.get('project');

    if (!project_name) {
        window.Utils.showToast(`No project is selected!`, "error");
        return;
    }

    const response = await fetch(`/api/logs/${project_name}/download-logs`, {
        method: "POST",
        body: JSON.stringify({
            "loglevels": checkcked_log_levels,
            "dates": dates
        })
    })
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "logs.xlsx";
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
})

document.addEventListener("DOMContentLoaded", async () => {
    set_checking_loglevels();
    await getProjects();
    await getLogs();
})