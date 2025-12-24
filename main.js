// @ts-check

async function main() {
	const cmd = new exec.Command('ls', '/');
	cmd.useStd(false, true, true);
	return await cmd.run();
}

main()
