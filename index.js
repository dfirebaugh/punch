import { editor } from "./editor.js";
import { renderJSON } from "./ast.js";
import { snippets } from "./examples.js";

document.addEventListener("DOMContentLoaded", () => {
  const resizer = document.getElementById("resizer");
  const editorContainer = document.querySelector(".editor-container");
  const outputContainer = document.querySelector(".output-container");

  let isResizing = false;
  let activeTab = localStorage.getItem("activeTab") || "js";
  let goInstance;
  let wasmInstance;
  let wasmRunning = false;
  let vimModeEnabled = localStorage.getItem("vimModeEnabled") === "true";

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

  async function initializeWasm() {
    if (wasmRunning) {
      console.warn("WASM instance already running");
      return;
    }
    goInstance = new Go();
    const result = await WebAssembly.instantiateStreaming(
      fetch("main.wasm"),
      goInstance.importObject,
    );
    wasmInstance = result.instance;
    goInstance.run(wasmInstance);
    wasmRunning = true;
  }

  async function ensureWasmRunning() {
    if (!wasmRunning) {
      await initializeWasm();
    }
  }

  function handleWasmError(error) {
    console.error("WASM error:", error);
    wasmRunning = false;
    const outputElement = document.getElementById("output");
    outputElement.innerText = "An error occurred, failed to generate";
    initializeWasm()
      .then(() => {
        console.log("WASM instance restarted");
      })
      .catch((initError) => {
        console.error("Failed to restart WASM instance:", initError);
      });
  }

  function parseCode() {
    ensureWasmRunning()
      .then(() => {
        try {
          switchTab("ast");
          const source = editor.getValue().trim();
          const ast = parse(source);
          const outputElement = document.getElementById("output");
          outputElement.innerHTML = "";
          renderJSON(JSON.parse(ast), outputElement);
        } catch (error) {
          handleWasmError(error);
        }
      })
      .catch((error) => {
        console.error("Failed to parse code:", error);
      });
  }

  document.getElementById("tab-ast").addEventListener("click", () => {
    activeTab = "ast";
    localStorage.setItem("activeTab", activeTab);
    parseCode();
  });

  document.getElementById("tab-lex").addEventListener("click", () => {
    activeTab = "lex";
    localStorage.setItem("activeTab", activeTab);
    ensureWasmRunning()
      .then(() => {
        try {
          switchTab("lex");
          const source = editor.getValue().trim();
          const tokens = lex(source);
          const outputElement = document.getElementById("output");
          outputElement.innerHTML = `<pre>${JSON.parse(tokens).join("\n")}</pre>`;
        } catch (error) {
          handleWasmError(error);
        }
      })
      .catch((error) => {
        console.error("Failed to lex code:", error);
      });
  });

  document.getElementById("tab-wat").addEventListener("click", () => {
    activeTab = "wat";
    localStorage.setItem("activeTab", activeTab);
    ensureWasmRunning()
      .then(() => {
        try {
          switchTab("wat");
          const source = editor.getValue().trim();
          const watCode = generateWAT(source);
          const outputElement = document.getElementById("output");
          outputElement.innerHTML = "";
          const watEditor = CodeMirror(outputElement, {
            value: watCode,
            mode: "wat",
            lineNumbers: true,
            theme: "dracula",
            readOnly: true,
          });
        } catch (error) {
          handleWasmError(error);
        }
      })
      .catch((error) => {
        console.error("Failed to generate WAT code:", error);
      });
  });

  document.getElementById("tab-js").addEventListener("click", () => {
    activeTab = "js";
    localStorage.setItem("activeTab", activeTab);
    ensureWasmRunning()
      .then(() => {
        try {
          switchTab("js");
          const source = editor.getValue().trim();
          const jsCode = generateJS(source);
          const formattedJsCode = prettier.format(jsCode, {
            parser: "babel",
            plugins: prettierPlugins,
          });
          const outputElement = document.getElementById("output");
          outputElement.classList.add("CodeMirror-js");
          outputElement.innerHTML = "";
          const jsEditor = CodeMirror(outputElement, {
            value: formattedJsCode,
            mode: "javascript",
            lineNumbers: true,
            theme: "dracula",
            readOnly: true,
          });

          document.getElementById("js-console").style.display = "block";
        } catch (error) {
          handleWasmError(error);
        }
      })
      .catch((error) => {
        console.error("Failed to generate JS code:", error);
      });
  });

  function switchTab(tab) {
    document
      .querySelectorAll(".tab")
      .forEach((t) => t.classList.remove("active"));
    document.getElementById(`tab-${tab}`).classList.add("active");

    const outputElement = document.getElementById("output");
    if (tab !== "ast") {
      outputElement.innerHTML = "";
    }

    if (tab !== "js") {
      document.getElementById("js-console").style.display = "none";
    }
  }

  function fetchAndRenderAST() {
    ensureWasmRunning()
      .then(() => {
        try {
          const source = editor.getValue().trim();

          switch (activeTab) {
            case "ast":
              const ast = parse(source);
              const outputElement = document.getElementById("output");
              outputElement.innerHTML = "";
              renderJSON(JSON.parse(ast), outputElement);
              break;
            case "lex":
              const tokens = lex(source);
              const outputElementLex = document.getElementById("output");
              outputElementLex.innerHTML = `<pre>${tokens}</pre>`;
              break;
            case "wat":
              //const watCode = generateWAT(source);
              //const outputElementWat = document.getElementById("output");
              //outputElementWat.innerHTML = "";
              //const watEditor = CodeMirror(outputElementWat, {
              //  value: watCode,
              //  mode: "wat",
              //  lineNumbers: true,
              //  theme: "dracula",
              //  readOnly: true,
              //});
              break;
            case "js":
              const jsCode = generateJS(source);
              const formattedJsCode = prettier.format(jsCode, {
                parser: "babel",
                plugins: prettierPlugins,
              });
              const outputElementJs = document.getElementById("output");
              outputElementJs.innerHTML = "";
              const jsEditor = CodeMirror(outputElementJs, {
                value: formattedJsCode,
                mode: "javascript",
                lineNumbers: true,
                theme: "dracula",
                readOnly: true,
              });

              document.getElementById("js-console").style.display = "block";
              break;
            default:
              console.error("Unknown tab:", activeTab);
          }
        } catch (error) {
          handleWasmError(error);
        }
      })
      .catch((error) => {
        console.error("Failed to fetch and render AST:", error);
      });
  }

  // Listen for :w command in Vim mode
  CodeMirror.Vim.defineEx("write", "w", () => {
    fetchAndRenderAST();
  });

  initializeWasm()
    .then(() => {
      document.getElementById(`tab-${activeTab}`).click();
    })
    .catch((error) => {
      console.error("Failed to initialize WASM:", error);
    });

  document.getElementById("run-js").addEventListener("click", () => {
    const outputElement = document.getElementById("output");
    const jsCode = outputElement
      .querySelector(".CodeMirror")
      .CodeMirror.getValue();
    const consoleOutput = document.getElementById("console-output");

    consoleOutput.textContent = "";

    const originalConsoleLog = console.log;
    const originalConsoleError = console.error;
    console.log = (...args) => {
      originalConsoleLog(...args);
      consoleOutput.textContent +=
        args
          .map((arg) =>
            typeof arg === "object" ? JSON.stringify(arg, null, 2) : arg,
          )
          .join(" ") + "\n";
    };
    console.error = (...args) => {
      originalConsoleError(...args);
      consoleOutput.textContent += `Error: ${args.map((arg) => (typeof arg === "object" ? JSON.stringify(arg, null, 2) : arg)).join(" ")}\n`;
    };

    try {
      const blob = new Blob([jsCode], { type: "application/javascript" });
      const url = URL.createObjectURL(blob);

      import(url)
        .then((module) => {
          URL.revokeObjectURL(url);
          console.log = originalConsoleLog;
          console.error = originalConsoleError;
        })
        .catch((error) => {
          console.error(`Error: ${error.message}`);
          URL.revokeObjectURL(url);
          console.log = originalConsoleLog;
          console.error = originalConsoleError;
        });
    } catch (error) {
      console.log = originalConsoleLog;
      console.error = originalConsoleError;
      console.error(`Error: ${error.message}`);
    }
  });

  document.getElementById("clear-console").addEventListener("click", () => {
    const consoleOutput = document.getElementById("console-output");
    consoleOutput.textContent = "";
  });

  document.getElementById("compile-code").addEventListener("click", () => {
    fetchAndRenderAST();
  });

  document.getElementById("toggle-vim").addEventListener("click", () => {
    vimModeEnabled = !vimModeEnabled;
    localStorage.setItem("vimModeEnabled", vimModeEnabled);
    editor.setOption("keyMap", vimModeEnabled ? "vim" : "default");
    document.getElementById("toggle-vim").textContent = vimModeEnabled
      ? "Disable Vim Mode"
      : "Enable Vim Mode";
  });

  editor.setOption("keyMap", vimModeEnabled ? "vim" : "default");
  document.getElementById("toggle-vim").textContent = vimModeEnabled
    ? "Disable Vim Mode"
    : "Enable Vim Mode";

  const snippetSelect = document.getElementById("code-snippets");
  Object.keys(snippets).forEach((key) => {
    const option = document.createElement("option");
    option.value = key;
    option.textContent = key.replace("snippet", "Snippet ");
    snippetSelect.appendChild(option);
  });

  snippetSelect.addEventListener("change", (event) => {
    const snippetKey = event.target.value;
    if (snippets[snippetKey]) {
      editor.setValue(snippets[snippetKey]);
    }
  });
});
