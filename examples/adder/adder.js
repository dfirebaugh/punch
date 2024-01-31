const fs = require('node:fs');

const wasmBuffer = fs.readFileSync('./adder.wasm');
WebAssembly.instantiate(wasmBuffer).then(wasmModule => {
  const { add_two, add_four } = wasmModule.instance.exports;
  const sum = add_four(1,1,1+2-2,1);
  const sum2 = add_two(2, 20);
  console.log(sum);
  console.log(sum2);
});

