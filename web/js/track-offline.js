
const OFFLINE_STATE = (0 <= document.URL.indexOf("/offline"));

window.addEventListener("load", () => {
  if(window.navigator.onLine) goOnline();
  else goOffline();
});

function goOffline() {
  if(!OFFLINE_STATE)
  window.open("/offline", "_self");
}

function goOnline() {
  if(OFFLINE_STATE)
  window.open("/", "_self");
}

window.addEventListener("offline", goOffline);
window.addEventListener("online", goOnline);
