CodeMirror.defineSimpleMode("wat", {
  start: [
    {
      regex:
        /\b(module|func|param|result|local|get|set|call|if|then|else|loop|block|br|br_if|return|i32|i64|f32|f64|memory|export|import|data|global|mut|const|store|load|offset)\b/,
      token: "keyword",
    },
    { regex: /"(?:[^\\]|\\.)*?"/, token: "string" },
    { regex: /;.*$/, token: "comment" },
    { regex: /-?\d+(\.\d+)?/, token: "number" },
    { regex: /[+\-/*=<>!&|]/, token: "operator" },
    { regex: /[\[\]{}()]/, token: "bracket" },
    { regex: /\$[a-zA-Z_][a-zA-Z0-9_]*/, token: "variable" },
  ],
  meta: {
    lineComment: ";",
  },
});
