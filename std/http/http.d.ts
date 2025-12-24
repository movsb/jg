declare global {
	namespace http {
		export class Response {
			text(): Promise<string>;
			blob(): Promise<ArrayBuffer>;
		}
		function get(url: string): Promise<Response>;
	}
}
export {};
