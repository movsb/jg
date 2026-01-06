async function main() {
	const stat = fs.stat('main.js');
	console.log(stat.name, stat.size, stat.modTime.unix(), stat.isDir, stat.isRegular);
}

main();
