const HOST = "com.gputoggle.helper";

const toggle = document.getElementById("toggle");
const stateLabel = document.getElementById("state-label");
const banner = document.getElementById("banner");
const btnRelaunch = document.getElementById("btn-relaunch");
const btnLater = document.getElementById("btn-later");
const errorDiv = document.getElementById("error");

// Desired state set by the user — sent to the helper on relaunch.
let pendingEnabled = null;

function showError(msg) {
  errorDiv.textContent = msg;
  errorDiv.style.display = "block";
  toggle.disabled = true;
}

function renderState(enabled) {
  toggle.checked = enabled;
  toggle.disabled = false;
  stateLabel.textContent = enabled ? "Enabled" : "Disabled";
  errorDiv.style.display = "none";
}

function send(msg, cb) {
  chrome.runtime.sendNativeMessage(HOST, msg, (response) => {
    if (chrome.runtime.lastError) {
      cb(null, chrome.runtime.lastError.message);
      return;
    }
    cb(response, null);
  });
}

// Load initial state.
send({ type: "getState" }, (resp, err) => {
  if (err) {
    showError("Helper not connected. Run install.ps1 to set it up.\n\n" + err);
    return;
  }
  if (resp.type === "error") {
    showError(resp.message);
    return;
  }
  renderState(resp.enabled);
});

// Toggle change — just update UI and show relaunch banner.
// The actual write happens in the helper after Chrome exits (see restart handler).
toggle.addEventListener("change", () => {
  pendingEnabled = toggle.checked;
  stateLabel.textContent = toggle.checked ? "Enabled (pending relaunch)" : "Disabled (pending relaunch)";
  banner.style.display = "block";
});

// Relaunch now — send desired state with the restart request.
btnRelaunch.addEventListener("click", () => {
  if (pendingEnabled === null) return;
  btnRelaunch.disabled = true;
  send({ type: "restart", enabled: pendingEnabled }, () => {
    // Chrome will close this popup when it relaunches.
  });
});

// Later — hide the banner; user can close Chrome manually to apply.
btnLater.addEventListener("click", () => {
  banner.style.display = "none";
});
