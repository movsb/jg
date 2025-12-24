// @ts-check

async function main() {
	const rsp = await http.get('http://baidu.com');
	const text = await rsp.text();
	console.log('length:', text);
}

main()
