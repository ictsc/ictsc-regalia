import "./index.css";
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
