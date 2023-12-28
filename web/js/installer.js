
const IS_INSTALLED = (0 <= document.URL.indexOf("?installed"));

var awaitsPrompt = !IS_INSTALLED;
var awaitedPrompt = undefined;
window.addEventListener("beforeinstallprompt", (ev) => {
  awaitedPrompt = ev;
});

const INSTALL_PROMPT_CALLER = async () => {
  if(awaitsPrompt && awaitedPrompt) {
    awaitedPrompt.prompt();
    awaitsPrompt = false;
  }
};
if(awaitsPrompt) {
  window.addEventListener("click", INSTALL_PROMPT_CALLER);
  window.addEventListener("keydown", INSTALL_PROMPT_CALLER);
}

if(IS_INSTALLED) {
  const PRE_INSTALL_WINDOW_OPEN = window.open;
  window.open = (url, target, features) => {
    return PRE_INSTALL_WINDOW_OPEN(url + "?installed", target, features);
  };
}
