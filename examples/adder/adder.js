const fs = require('node:fs');

const wasmBuffer = fs.readFileSync('./adder.wasm');
WebAssembly.instantiate(wasmBuffer).then(wasmModule => {
  const { add_two, add_four } = wasmModule.instance.exports;
  const sum = add_four(5, 6, 2, 4);
  const sum2 = add_two(5, 6);
  console.log(sum);
  console.log(sum2);
});

