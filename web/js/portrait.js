
var portraitMode;
function eveluatePortraitMode() {
  portraitMode = window.innerWidth < window.innerHeight;
  if(portraitMode) {
    if(!document.body.classList.contains("portrait"))
    document.body.classList.add("portrait");
  } else {
    if(document.body.classList.contains("portrait"))
    document.body.classList.remove("portrait");
  }
}

window.addEventListener("load", eveluatePortraitMode);
window.addEventListener("resize", eveluatePortraitMode);
