// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
const vscode = require('vscode');
const cp = require('child_process');
const fs = require('fs');

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
		// The code you place here will be executed every time your command is executed

		// Display a message box to the user
		vscode.window.showInformationMessage('Hello World from yaml-readme!');

		cp.exec('pwd', (err, stdout, stderr) => {
			console.log('stdout: ' + stdout);
			console.log('stderr: ' + stderr);
			if (err) {
				console.log('error: ' + err);
			}
		});

		if(vscode.workspace.workspaceFolders !== undefined) {
			let wf = vscode.workspace.workspaceFolders[0].uri.path ;
			let f = vscode.workspace.workspaceFolders[0].uri.fsPath ; 
		
			// let message = `YOUR-EXTENSION: folder: ${wf} - ${f}` ;
		
			// vscode.window.showInformationMessage(message);

			let filename = vscode.window.activeTextEditor.document.fileName
			console.log(filename)
			vscode.window.showInformationMessage(filename);
			let metadata = getFirstLine(filename)
			vscode.window.showInformationMessage(metadata + "===" + metadata.startsWith("#!yaml-readme"));
			if (metadata.startsWith("#!yaml-readme")) {
				metadata = metadata.replace("#!yaml-readme ", "")
				vscode.window.showInformationMessage(metadata);

				let pattern = ""
				let output = ""
				const items = metadata.split(" ")
				for (var i = 0; i < items.length; i++) {
					const item = items[i]
					if (item == "-p") {
						pattern = wf + "/" + items[++i]
					} else if (item == "--output") {
						output = wf + "/" + items[++i]
					}
				}

				vscode.window.showInformationMessage(`yaml-readme -p "${pattern}" -t "${filename}" > ${output}`)
				cp.exec(`yaml-readme -p "${pattern}" -t "${filename}" > ${output}`, (err, stdout, stderr) => {
					console.log('stdout: ' + stdout);
					console.log('stderr: ' + stderr);
					if (err) {
						console.log('error: ' + err);
					}
					vscode.commands.executeCommand("markdown.showPreviewToSide", vscode.Uri.file(`${output}`));
				});
			}
		}  else {
			let message = "YOUR-EXTENSION: Working folder not found, open a folder an try again" ;
		
			vscode.window.showErrorMessage(message);
		}

		// initMarkdownPreview(context);
	});

	context.subscriptions.push(disposable);
}

async function initMarkdownPreview(context) {
    const panel = vscode.window.createWebviewPanel(
        // Webview id
        'liveHTMLPreviewer',
        // Webview title
        '[Preview]',
        // This will open the second column for preview inside editor
        2,
        {
            // Enable scripts in the webview
            enableScripts: true,
            retainContextWhenHidden: true,
        }
    );

	
	let wf = vscode.workspace.workspaceFolders[0].uri.path ;
	let f = vscode.workspace.workspaceFolders[0].uri.fsPath ; 

	let message = `YOUR-EXTENSION: folder: ${wf} - ${f}` ;

	vscode.window.showInformationMessage(message);

	cp.exec(`yaml-readme -p "${wf}/what/items/job-*.yaml" -t "${wf}/what/jobs.tpl"`, (err, stdout, stderr) => {
		console.log('stdout: ' + stdout);
		console.log('stderr: ' + stderr);
		panel.webview.html = stdout
		if (err) {
			console.log('error: ' + err);
		}

		vscode.commands.executeCommand("markdown.showPreview", stdout);
	});
}

// this method is called when your extension is deactivated
function deactivate() {}

module.exports = {
	activate,
	deactivate
}
