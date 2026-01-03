async function main() {
	const name = 'vim';
	const cmd = $`${name}`;
	cmd.useStd(true, true, true);
	return await cmd.run();
}

main();
