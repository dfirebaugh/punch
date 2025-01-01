import { renderJSON } from "/ast.js";

const editor = CodeMirror(document.getElementById('editor'), {
  mode: "rust",
  lineNumbers: true,
  theme: "dracula",
  keyMap: "vim",
  value: `
pkg main

bool is_eq(i32 a, i32 b) {
    return a == b
}

pub i32 add_two(i32 x, i32 y, i32 z) {
    println("x = {}, y = {}", x, y, z)
    println("Hello, World!")
    return x + y
}

pub i32 add_four(i32 a, i32 b, i32 c, i32 d) {
    return a + b + c + d
}

      `.trim()
});

const highlightCode = (startLine, startCol, endLine, endCol) => {
  editor.getAllMarks().forEach(mark => mark.clear());
  const from = { line: startLine - 1, ch: startCol - 1 };
  const to = { line: endLine - 1, ch: endCol - 1 };
  editor.markText(from, to, { className: 'highlighted' });
};

let lastSource = "";

export const parseCode = async () => {
  const source = editor.getValue().trim();

  if (source === lastSource) {
    return;
  }

  lastSource = source;

  try {
    const response = await fetch('/parse', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ source }),
    });

    if (!response.ok) {
      throw new Error(`Server error: ${response.statusText}`);
    }

    const ast = await response.json();

    const outputElement = document.getElementById('output');
    outputElement.innerHTML = '';

    renderJSON(ast, outputElement);
  } catch (error) {
    const outputElement = document.getElementById('output');
    outputElement.innerText = `Error: ${error.message}`;
  }
};

editor.on('change', parseCode);

