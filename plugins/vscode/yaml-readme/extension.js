// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
const vscode = require('vscode');
const cp = require('child_process');
const fs = require('fs');
const cmd = require('./command');
const {TestTask, HelloReply} = require('./server_pb');
const {RunnerClient} = require('./server_grpc_web_pb');
global.XMLHttpRequest = require('xhr2');


const grpc = require('@grpc/grpc-js');

const protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = __dirname +'/server.proto';

const packageDefinition = protoLoader.loadSync(
  PROTO_PATH, { 
   keepCase: true,
   longs: String,
   enums: String,
   defaults: true,
   oneofs: true,
});

const serverProto = grpc.loadPackageDefinition(packageDefinition).server;

const apiConsole = vscode.window.createOutputChannel("API Testing")

function getFirstLine(filePath) {
	const data = fs.readFileSync(filePath);
  return data.toString().split('\n')[0];
}

// this method is called when your extension is activated
// your extension is activated the very first time the command is executed

/**
 * @param {vscode.ExtensionContext} context
 */
function activate(context) {

	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	console.log('Congratulations, your extension "yaml-readme" is now active!');

	// The command has been defined in the package.json file
	// Now provide the implementation of the command with  registerCommand
	// The commandId parameter must match the command field in package.json
	let disposable = vscode.commands.registerCommand('yaml-readme.helloWorld', function () {
		if(vscode.workspace.workspaceFolders !== undefined) {
			let wf = vscode.workspace.workspaceFolders[0].uri.path ;

			let filename = vscode.window.activeTextEditor.document.fileName
			let metadata = getFirstLine(filename)
			// vscode.window.showInformationMessage(metadata + "===" + metadata.startsWith("#!yaml-readme"));
			if (metadata.startsWith("#!yaml-readme")) {
				let command = cmd.generateCommand(metadata, wf, filename)

				// vscode.window.showInformationMessage(`yaml-readme -p "${pattern}" -t "${filename}" > ${output}`)
				cp.exec(`${command[0]} > ${command[1]}`, (err) => {
					if (err) {
						console.log('error: ' + err);
					}
					vscode.commands.executeCommand("markdown.showPreviewToSide", vscode.Uri.file(`${command[1]}`));
				});
			}
		}  else {
			let message = "YOUR-EXTENSION: Working folder not found, open a folder an try again" ;
		
			vscode.window.showErrorMessage(message);
		}
	});
	let atest = vscode.commands.registerCommand('atest', function() {
		if(vscode.workspace.workspaceFolders !== undefined) {
			let filename = vscode.window.activeTextEditor.document.fileName
			const addr = vscode.workspace.getConfiguration().get('yaml-readme.server')
			apiConsole.show()

			var c = new RunnerClient(addr)
			let task = new TestTask()
			task.setKind('suite')

			let editor = vscode.window.activeTextEditor
			if (editor) {
				let selection = editor.selection
				let text = editor.document.getText(selection)
				if (text !== undefined && text !== '') {
					task.setData(text)
				}
			}

			if (task.getData() === undefined || task.getData() === '') {
				const data = fs.readFileSync(filename);
				task.setData(data.toString())
			}

			apiConsole.appendLine("start to run")
			try {
				let a = c.run(task, {}, (err,resp) => {
					console.log(err,resp)
				})
				console.log(a)
			} catch(e){
				console.log(e)
			}

			const client = new serverProto.Runner(addr, grpc.credentials.createInsecure());
			let employeeIdList = [1,2,3];
			client.run({
				kind: "suite",
				data: data.toString()
			} , function(err, response) {
			  console.log('Data:', response); // API response
			  console.log(err);
			 });
		}  else {
			let message = "YOUR-EXTENSION: Working folder not found, open a folder an try again" ;
		
			vscode.window.showErrorMessage(message);
		}
	})

	context.subscriptions.push(disposable,atest);
}

// this method is called when your extension is deactivated
function deactivate() {}

module.exports = {
	activate,
	deactivate
}
