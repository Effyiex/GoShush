
self.addEventListener("install", (ev) => {
  ev.waitUntil(caches.open("sw-cache").then((cache) => {
    return cache.addAll([
      "/fonts/Nunito.ttf",
      "/fonts/Acme.ttf",
      "/offline",
      "/js/install.js",
      "/js/offline.js",
      "/css/master.css"
    ]);
  }));
});

self.addEventListener("fetch", (ev) => {
  ev.respondWith(caches.match(ev.request).then((response) => {
    return response || fetch(ev.request);
  }));
});
