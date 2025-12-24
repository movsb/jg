declare global {
	namespace io {
		export class Reader {}
	}
	namespace fs {
		function saveToFile(filePath: string, data: io.Reader): Promise<number>;
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
