import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";

if ("serviceWorker" in navigator) {
  navigator.serviceWorker.register("/sw.js").then((registration) => {
    console.log("SW registered: ", registration);
  });
}

// request notification permission

if ("Notification" in window) {
  Notification.requestPermission().then((permission) => {
    console.log("Notification permission:", permission);
  });
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
