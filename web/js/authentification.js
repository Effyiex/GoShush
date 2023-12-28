
var user = JSON.parse(localStorage.getItem("user"));

async function sha256(message) {
  const msgBuffer = new TextEncoder().encode(message); 
  const hashBuffer = await crypto.subtle.digest('SHA-256', msgBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));            
  const hashHex = hashArray.map(b => b.toString(16).padStart(2, '0')).join("");
  return hashHex;
}

window.addEventListener("load", () => {

  if(!document.querySelector("#login")) {
    if(!user) window.open("/login", "_self");
    return;
  }

  const loginButton = document.querySelector("#login #submit");
  const authData = document.querySelector("#login #authdata");
  const registerToggle = authData.querySelector("#register"); 
  const keyInput = authData.querySelector("#key");
  keyInput.style.display = "none";

  authData.querySelectorAll("input").forEach(input => {
    input.addEventListener("keyup", ev => {
      if(ev.key == "Enter") 
      loginButton.click();
    });
  });

  loginButton.addEventListener("click", async () => {

    var authPacket = {
      username: authData.querySelector("#username").value,
      password: await sha256(authData.querySelector("#password").value),
      key: ""
    };

    if(keyInput.value.length > 0) 
    authPacket.key = keyInput.value;

    fetch("api/auth", {
      method: "POST",
      body: JSON.stringify(authPacket)
    }).then(data => data.json()).then(user => {
      if(user != null && user.name.length > 0) {
        console.log(user);
        localStorage.setItem("user", JSON.stringify(user));
        window.open("/frame", "_self");
      } else {
        alert("Login failed");
      }
    });

  });

  registerToggle.addEventListener("click", () => {
    if(registerToggle.checked) keyInput.style.display = "block";
    else keyInput.style.display = "none";
    keyInput.value = "";
  });

});
