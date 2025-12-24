declare global {
	namespace runtime {
		/**
		 * Current OS name (from Go).
		 */
		export const os: string;
		/**
		 * Current ARCH name (from Go).
		 */
		export const arch: string;
		/**
		 * Program arguments passed from command line.
		 */
		export const args: string[];
	}
	namespace io {
		export class Reader {}
	}
	namespace fs {
		function saveToFile(filePath: string, data: io.Reader): Promise<number>;
		/**
		 * Checks if the specified file or directory exists.
		 * 
		 * @param filePath 
		 */
		function fileExists(filePath: string): boolean;
		/**
		 * Checks if the specified file or directory exists.
		 * 
		 * @param filePath 
		 * @param types     File types to match.
		 * 
		 * e.g: 'fd' means either file or directory exists. 'fx' means file exists and executable.
		 * 
		 *  OR-ed
		 *      - 'f' file
		 *      - 'd' directory
		 *      - 'l' soft link
		 *      - 's' socket
		 * 
		 *  AND-ed
		 *      - 'x' executable
		 *      - 'w' writable
		 *      - 'r' readable
		 */
		function fileExists(filePath: string, types: string): boolean;
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
		function $(literals: TemplateStringsArray, ...interpolates: any[]): Command;
	}
	function $(literals: TemplateStringsArray, ...interpolates: any[]): exec.Command;
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
