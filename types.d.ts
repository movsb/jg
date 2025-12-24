declare global {
	namespace io {
		export class Reader {}
	}
	namespace fs {
		function saveToFile(filePath: string, data: io.Reader): Promise<number>;
	}
	namespace archive {
		export class TarReader {
			constructor(input: io.Reader);
			extractTo(dir: string): Promise<void>;
		}
	}
	namespace exec {
		export class Command {
			constructor(cmd: string, ...args: string[]);
			run(): Promise<void>;
			useStd(stdin: boolean, stdout: boolean, stderr: boolean);
		}
	}
	namespace http {
		export class Response {
			text(): Promise<string>;
			blob(): Promise<ArrayBuffer>;
			reader(): io.Reader;
		}
		function get(url: string): Promise<Response>;
	}
}
export {};
