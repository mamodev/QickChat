const cacheName = "my-cache";
const cacheVersion = "v1";
const cacheKey = `${cacheName}-${cacheVersion}`;

const filesToCache = ["/", "/index.js", "/index.css", "logo.png", "manifest.json", "favicon.ico"];

self.addEventListener("install", (event) => {
  const requests = filesToCache.map((file) => {
    return fetch(new Request(file, { cache: "reload" }));
  });

  event.waitUntil(
    Promise.all(requests).then((responses) => {
      return caches.open(cacheKey).then((cache) => {
        console.log("[oninstall] Cached offline page", responses);
        responses.forEach((response, index) => {
          cache.put(filesToCache[index], response);
        });
      });
    })
  );
});

self.addEventListener("activate", (event) => {
  console.log("service worker activated");
  event.waitUntil(
    caches.keys().then((cacheKeys) => {
      return Promise.all(
        cacheKeys.map((name) => {
          if (name !== cacheName) {
            return caches.delete(name);
          }
        })
      );
    })
  );
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    })
  );
});

self.addEventListener("message", (event) => {
  if (event.data && event.data.type === "checkForUpdate") {
    self.skipWaiting();
  }
});

self.addEventListener("push", (event) => {
  const options = {
    body: event.data.text(),
  };

  console.log("event push", event);

  event.waitUntil(self.registration.showNotification("Notification", options));
});
