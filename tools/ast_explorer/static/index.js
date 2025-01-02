import { editor, fetchAndRenderAST } from "/editor.js";

const resizer = document.getElementById("resizer");
const editorContainer = document.querySelector(".editor-container");
const outputContainer = document.querySelector(".output-container");

let isResizing = false;

resizer.addEventListener("mousedown", (e) => {
  isResizing = true;
  document.body.style.cursor = "col-resize";
});

document.addEventListener("mousemove", (e) => {
  if (!isResizing) return;

  const containerRect = document
    .querySelector(".container")
    .getBoundingClientRect();
  const newEditorWidth = e.clientX - containerRect.left;
  const newOutputWidth =
    containerRect.width - newEditorWidth - resizer.offsetWidth;

  if (newEditorWidth > 100 && newOutputWidth > 100) {
    editorContainer.style.width = `${newEditorWidth}px`;
    outputContainer.style.width = `${newOutputWidth}px`;
  }
});

document.addEventListener("mouseup", () => {
  isResizing = false;
  document.body.style.cursor = "default";
});

function parseCode() {
  switchTab("ast");
  fetchAndRenderAST();
}
document.getElementById("tab-ast").addEventListener("click", () => {
  parseCode();
});
parseCode();

document.getElementById("tab-lex").addEventListener("click", () => {
  switchTab("lex");
  const source = editor.getValue().trim();

  fetch("/lex", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ source }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Server error: ${response.statusText}`);
      }
      return response.json();
    })
    .then((tokens) => {
      const outputElement = document.getElementById("output");
      outputElement.innerHTML = `<pre>${JSON.stringify(tokens, null, 2)}</pre>`;
    })
    .catch((error) => {
      const outputElement = document.getElementById("output");
      outputElement.innerText = `Error: ${error.message}`;
    });
});

document.getElementById("tab-wat").addEventListener("click", () => {
  switchTab("wat");
  const source = editor.getValue().trim();

  fetch("/wat", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ source }),
  })
    .then((response) => {
      if (!response.ok) {
        throw new Error(`Server error: ${response.statusText}`);
      }
      return response.text();
    })
    .then((watCode) => {
      const outputElement = document.getElementById("output");
      outputElement.innerHTML = "";
      const watEditor = CodeMirror(outputElement, {
        value: watCode,
        mode: "wat",
        lineNumbers: true,
        theme: "dracula",
        readOnly: true,
      });
    })
    .catch((error) => {
      const outputElement = document.getElementById("output");
      outputElement.innerText = `Error: ${error.message}`;
    });
});

function switchTab(tab) {
  document
    .querySelectorAll(".tab")
    .forEach((t) => t.classList.remove("active"));
  document.getElementById(`tab-${tab}`).classList.add("active");

  const outputElement = document.getElementById("output");
  if (tab !== "ast") {
    outputElement.innerHTML = ""; // Clear output for lexed tokens and WAT code
  }
}

document.getElementById("tab-wat").click();
