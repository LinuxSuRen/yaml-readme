const assert = require('assert');

// You can import and use all API from the 'vscode' module
// as well as import your extension to test it
// const vscode = require('vscode');
const myExtension = require('../../command');

let cmd = myExtension.generateCommand("#!yaml-readme -p data/*.yaml --output README.md --group-by kind --sort-by kind","wf","filename")
assert.equal(cmd[0], "yaml-readme -t filename -p wf/data/*.yaml --group-by kind --sort-by kind")
assert.equal(cmd[1], "wf/README.md")
