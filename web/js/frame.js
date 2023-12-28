
var frame;

window.addEventListener("load", () => {

  document
    .querySelector("#usercontrol #logoff")
    .addEventListener("click", () => {
      localStorage.removeItem("user");
      window.open("/login", "_self");
  });

  document.querySelector("#usercontrol #welcome username").innerHTML = user.name;

  var savedFrameSrc = localStorage.getItem("framesrc");
  if(!savedFrameSrc) {
    savedFrameSrc = "frames/chats";
    localStorage.setItem("framesrc", savedFrameSrc);
  }

  frame = document.querySelector("#content #display");
  frame.src = savedFrameSrc;
  frame.addEventListener("load", () => {
    document.title = frame.contentDocument.title;
  });

  document.querySelectorAll(".fr-invoke").forEach((element, _) => {
    element.classList.forEach(clazz => {
      if(!clazz.startsWith("fr-set-")) return;
      var setter = clazz.substring(7);
      console.log("Registering Invoker: \"frames/" + setter + "\" on: ");
      console.log(element);
      element.addEventListener("click", () => {
        document
          .querySelectorAll(".fr-invoke")
          .forEach(e => e.classList.remove("fr-active"));
        element.classList.add("fr-active");
        frame.src = "frames/" + setter + "?cache=" + new Date().getTime();
        localStorage.setItem("framesrc", "frames/" + setter);
        console.log("Opening: frames/" + setter);
      });
    });
  });

});
