// @ts-check

async function main() {
	const rsp = await http.get('http://localhost:4637');
	const body = rsp.reader();
	const tr = new archive.TarReader(body);
	return tr.extractTo('/tmp');
}

main()
