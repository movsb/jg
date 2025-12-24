// @ts-check

async function main() {
	const rsp = await http.get('http://baidu.com');
	const body = rsp.reader();
	return await fs.saveToFile('baidu.com.txt', body);
}

main()
