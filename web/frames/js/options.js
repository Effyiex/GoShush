
const APP_SETTINGS = {
  
  sliderOptions: {
    invert: 10,
    hue: 0,
    zoom: 100
  },

  boolOptions: {
    rainbowHue: false
  },

  rainbowInterval: undefined,
  set rainbowDegree(value) {
    localStorage.setItem("rainbow-degree", value);
  },
  get rainbowDegree() {
    return parseInt(localStorage.getItem("rainbow-degree"));
  },
  toggleRainbowHue: function(toggleElement) {
    this.boolOptions.rainbowHue = !this.boolOptions.rainbowHue;
    this.save();
    if(!this.boolOptions.rainbowHue) {
      clearInterval(this.rainbowInterval);
      this.rainbowInterval = undefined;
      if(toggleElement) {
        toggleElement.classList.remove("toggled");
        console.log("Removed 'toggled'!");
      }
      this.apply();
    } else {
      if(toggleElement) {
        toggleElement.classList.add("toggled");
        console.log("Added 'toggled'!");
      }
      this.rainbowInterval = setInterval(() => {
        if(this.rainbowDegree < 360) this.rainbowDegree += 4;
        else this.rainbowDegree = 0;
        var filters = "hue-rotate(" + this.rainbowDegree + "deg)";
        document.body.style.filter.split(" ").forEach(f => {
          if(f.indexOf("hue-rotate") < 0)
          filters += f;
        });
        document.body.style.filter = filters;
      }, 64);
    }
  },

  load: function() {
    const sSliderOptions = localStorage.getItem("slider-options");
    if(sSliderOptions)
      this.sliderOptions = JSON.parse(sSliderOptions);
    const sBoolOptions = localStorage.getItem("bool-options");
    console.log("Loaded Bool Options: ")
    console.log(sBoolOptions);
    if(sBoolOptions)
      this.boolOptions = JSON.parse(sBoolOptions);
    console.log("Parsed Bool Options: ")
    console.log(this.boolOptions);
  },

  save: function() {
    localStorage.setItem("slider-options", JSON.stringify(this.sliderOptions));
    localStorage.setItem("bool-options", JSON.stringify(this.boolOptions));
  },

  apply: function() {
    if(this.sliderOptions.zoom)
    document.body.style.zoom = this.sliderOptions.zoom + "%";
    var filterList = [];
    if(
      this.sliderOptions.invert
      && this.sliderOptions.invert > 0
    )
      filterList.push("invert(" + this.sliderOptions.invert + "%)");
    if(
      this.sliderOptions.hue 
      && this.sliderOptions.hue > 0
    ) filterList.push("hue-rotate(" + this.sliderOptions.hue * 3.6 + "deg)");
    const filter = filterList.join(" ");
    console.log("Filter: " + filter);
    document.body.style.filter = filter;
    if(this.boolOptions.rainbowHue)
    this.toggleRainbowHue(document.querySelector("#rainbow"));
  }

}

window.addEventListener("load", () => {
  console.log("pre-load: ")
  console.log(APP_SETTINGS);
  APP_SETTINGS.load();
  console.log("post-load: ")
  console.log(APP_SETTINGS);
  APP_SETTINGS.apply();
  console.log("post-apply: ")
  console.log(APP_SETTINGS);
});

if(0 <= document.URL.indexOf("/options"))
window.addEventListener("load", () => {

  document.querySelectorAll(".opt-slider").forEach((element, _) => {
    const label = element.querySelector('p');
    const sliderTitle = label.innerHTML.split(':')[0].toLowerCase().trim();
    const slider = element.querySelector(".slider");
    var value = APP_SETTINGS.sliderOptions[sliderTitle];
    if(value) 
      slider.value = value;
    else
      value = slider.value;
    label.innerHTML = sliderTitle[0].toUpperCase() + sliderTitle.substring(1) + ": " + value + '%';
    slider.addEventListener("input", () => {
      label.innerHTML = sliderTitle[0].toUpperCase() + sliderTitle.substring(1) + ": " + slider.value + '%';
      APP_SETTINGS.sliderOptions[sliderTitle] = slider.value;
      APP_SETTINGS.save();
      APP_SETTINGS.apply();
    });
  });

  document.querySelector("#rainbow").addEventListener(
    "click", 
    () => APP_SETTINGS.toggleRainbowHue(document.querySelector("#rainbow"))
  );

});
