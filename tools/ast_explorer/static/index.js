import { parseCode } from "/editor.js"

const resizer = document.getElementById('resizer');
const editorContainer = document.querySelector('.editor-container');
const outputContainer = document.querySelector('.output-container');

let isResizing = false;

resizer.addEventListener('mousedown', (e) => {
  isResizing = true;
  document.body.style.cursor = 'col-resize';
});

document.addEventListener('mousemove', (e) => {
  if (!isResizing) return;

  const containerRect = document.querySelector('.container').getBoundingClientRect();
  const newEditorWidth = e.clientX - containerRect.left;
  const newOutputWidth = containerRect.width - newEditorWidth - resizer.offsetWidth;

  if (newEditorWidth > 100 && newOutputWidth > 100) {
    editorContainer.style.width = `${newEditorWidth}px`;
    outputContainer.style.width = `${newOutputWidth}px`;
  }
});

document.addEventListener('mouseup', () => {
  isResizing = false;
  document.body.style.cursor = 'default';
});


parseCode();

