/// <reference types="react/canary" />
/// <reference types="react-dom/canary" />
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

const App = () => {
  return <div>aaa</div>;
};

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
