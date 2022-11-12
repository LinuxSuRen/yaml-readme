function generateCommand(metadata, wf, filename) {
	metadata = metadata.replace("#!yaml-readme ", "")

	let commands = ["yaml-readme", "-t", filename]
	let output = ""
	const items = metadata.split(" ")

	for (var i = 0; i < items.length; i++) {
		const item = items[i]
		if (item == "-p") {
			commands.push("-p", wf + "/" + items[++i])
		} else if (item == "--output") {
			output = wf + "/" + items[++i]
		} else if (item == "--group-by") {
			commands.push("--group-by", items[++i])
		} else if (item == "--sort-by") {
			commands.push("--sort-by", items[++i])
		}
	}

	return [commands.join(" "), output]
}

module.exports = {
	generateCommand
}
