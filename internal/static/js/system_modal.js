const core_settings = document.getElementById("core-settings");
const core_settings_modal_body = document.getElementById("system-modal-body");

async function getSystemCoreSettings() {
    try {
        const res = await fetch(`/api/core/settings`)
        const data = await res.json();
        if(res.ok){
            window.Utils.showToast("Core settings are fetched!", "info");
            renderModal(data);
        }
        else{
            window.Utils.showToast("Failed to fetch core settings!", "error");
        }
    } catch (error) {
        console.log(error);
        window.Utils.showToast("Failed to fetch core settings!", "error");
    }
}

function renderModal (data){
    console.log(data)
    core_settings_modal_body.innerHTML = '';

    core_settings_modal_body.innerHTML += `
        <div class="grid gap-4 mb-4 grid-cols-2">
            <div class="col-span-2">
                <label for="autolog_delete_days" class="block mb-2 text-sm font-medium text-gray-900">Auto Log Delete Days</label>
                <input type="number" value='${data?.autolog_delete_days ?? 0}' name="autolog_delete_days" id="autolog_delete_days" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5">
          </div>
        </div>
    `

    core_settings_modal_body.innerHTML += `
        <button type="submit" onclick="updateCoreSettings()" class="text-white inline-flex items-center bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center">
          Update Settings
        </button>
    `
}

async function updateCoreSettings() {
    const autolog_delete_days = document.getElementById("autolog_delete_days").value;

    if(!autolog_delete_days){
        window.Utils.showToast("Some input values are missing!", "error");
        return
    }

    try {
        const res = await fetch(`/api/core/settings`,{
            method:"POST",
            body:JSON.stringify({
                autolog_delete_days:parseInt(autolog_delete_days)
            })
        })
        const data = await res.json();
        if(res.ok){
            window.Utils.showToast("Core settings are updated!", "success");
            renderModal(data);
        }
        else{
            window.Utils.showToast("Failed to update core settings!", "error");
        }
    } catch (error) {
        window.Utils.showToast("Failed to update core settings!", "error");
    }
}

core_settings.addEventListener("click",async(e)=>{
    document.getElementById("core-system-modal-open").click();
})

document.addEventListener("DOMContentLoaded",async()=>{
    await getSystemCoreSettings();
})