
window.addEventListener("load", () => {
  document.querySelectorAll("#username").forEach((element) => {
    element.innerHTML = element.innerHTML.replace("%username%", user.name);
  });
  document.querySelector("#name-edit").addEventListener("click", async () => {
    while(true) {
      const promptIn = prompt("Your new username: ");
      if(promptIn != null) {
        const response = await fetch("../api/rename", {
          method: "POST",
          body: JSON.stringify({
            user: user.id,
            session: user.session,
            new_name: promptIn
          })
        });
        const status = await response.text();
        if(status == "success") break;
        else alert(status);
      } else break;
    }
  });
});
