async function main() {
	console.log(fs.fileExists('/', 'd'))
	console.log(fs.fileExists('/', 'f'))
	console.log(fs.fileExists('/etc/hosts', 'd'))
	console.log(fs.fileExists('/etc/hosts', 'f'))
	console.log(fs.fileExists('/etc/hosts', 'fx'))
	console.log(fs.fileExists('/bin/ls', 'fx'))
}

main()
