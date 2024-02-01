const fs = require('node:fs');

const wasmBuffer = fs.readFileSync('./adder.wasm');

const encode = function stringToIntegerArray(string, array) {
  const alphabet = "abcdefghijklmnopqrstuvwxyz";
  for (let i = 0; i < string.length; i++) {
    array[i] = alphabet.indexOf(string[i]);
  }
};

const decode = function (memory, offset) {
  const alphabet = "abcdefghijklmnopqrstuvwxyz";
  let string = "";
  let char = memory[offset];
  while (char !== 0) {
    string += alphabet[char];
    offset++;
    char = memory[offset];
  }
  return string;
};

const importObject = {
  imports: {
    println(offset) {
      const memoryBuffer = new Uint8Array(importObject.env.memory.buffer);
      console.log(decode(memoryBuffer, offset));
    },
  },
  env: {
    memory: new WebAssembly.Memory({ initial: 1 }) // 1 page = 64KiB
  }
};

WebAssembly.instantiate(wasmBuffer, importObject).then(wasmModule => {
  importObject.env.memory = wasmModule.instance.exports.memory;

  const { add_two, add_four } = wasmModule.instance.exports;

  const sum = add_four(1, 1, 2, 1);
  const sum2 = add_two(2, 20);

  console.log(`Sum from add_four: ${sum}`);
  console.log(`Sum from add_two: ${sum2}`);
});
