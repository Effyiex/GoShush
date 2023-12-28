
const CHAT_CFG = {
  id: undefined, // default: broadcast
  get label() {
    return (this.id ? this.id /* TODO: Fetch Chat Label */ : "Broadcast");
  },
  get bufferStorageKey() {
    return ("chat-" + (this.id ? this.id : "broadcast") + "-messages");
  }
};

const AUDIO_CONTEXT = new AudioContext();
const isMicrophoneAllowed = callback => {
  navigator.permissions.query({
    name: "microphone"
  }).then(status => callback(status.state === "granted"));
}

function queryChatMessages() {

  var latestMessageBuffer = JSON.parse(localStorage.getItem(CHAT_CFG.bufferStorageKey));
  if(!latestMessageBuffer) 
    latestMessageBuffer = [];

  if(CHAT_CFG.id) {
    
  } else { // Broadcast
    fetch(
      "../api/broadcast/query"
      + "?user=" + user.id
      + "&session=" + user.session, {
      method: "GET"
    })
      .then((response) => response.json())
      .then((messages) => {
        for(var i = 0; i < messages.length; i++)
        latestMessageBuffer.push(messages[i]);
      });
  }

  localStorage.setItem(CHAT_CFG.bufferStorageKey, JSON.stringify(latestMessageBuffer));

}

class ChatMessageContent {

  constructor(body, type) {
    this.body = body;
    this.type = type ? type : "text";
  }

}

function sendChatMessage(contents) {

  if(!Array.isArray(contents))
  contents = [contents];

  if(CHAT_CFG.id) {

  } else { // Broadcast
    fetch("../api/broadcast/send", {
      method: "POST",
      body: JSON.stringify({
        content: contents,
        sender: user.id,
        session: user.session
      })
    })
      .then((response) => response.text())
      .then((status) => {
        if(status == "success") queryChatMessages();
        else alert("Failed to send message.");
      });
  }

}

var memoRecorder;
const MEMO_CLICK_LISTENER = async (ev) => {

  isMicrophoneAllowed(allowed => {
    if(allowed && ev.target.classList.contains("disabled"))
    ev.target.classList.remove("disabled");
  });

  if(memoRecorder) {

    if(ev.target.classList.contains("recording"))
    ev.target.classList.remove("recording");

    memoRecorder.stop();
    memoRecorder = undefined;

  } else navigator.mediaDevices.getUserMedia({
    audio: true
  }).then((stream) => {

    if(!ev.target.classList.contains("recording"))
    ev.target.classList.add("recording");

    memoRecorder = new MediaRecorder(stream);

    let buffer = [];
    memoRecorder.ondataavailable = (event) => {
      buffer.push(event.data);
    };
    memoRecorder.onstop = async () => {
      const blob = new Blob(buffer, {
        type: "audio/ogg; codecs=opus"
      });
      sendChatMessage(new ChatMessageContent(Array.from(await blob.arrayBuffer()), "memo"));
      console.log(buffer);
      buffer = [];
    };
    memoRecorder.start();

  }).catch((err) => {
    console.log(err);
  });

};

window.addEventListener("load", () => {

  const urlParams = new URLSearchParams(window.location.search);
  CHAT_CFG.id = urlParams.get("id");

  queryChatMessages();

  document.querySelectorAll("*").forEach(element => {
    var child = element.firstChild;
    while(child) {
      if(child.nodeType === Node.TEXT_NODE) 
        child.textContent = child.textContent.replace("%chatlabel%", CHAT_CFG.label);
      child = child.nextSibling;
    }
  });

  const inputField = document.querySelector(".chat-text");

  const memo = document.querySelector(".chat-memo");
  isMicrophoneAllowed(allowed => {
    if(!allowed)
    memo.classList.add("disabled");
  });

  const send = document.querySelector(".chat-send");
  send.addEventListener("click", async () => {
    sendChatMessage(new ChatMessageContent(inputField.value));
    inputField.value = "";
    queryChatMessages();
  });

  memo.addEventListener("click", MEMO_CLICK_LISTENER);

});
