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
		/**
		 * Print panic info and exit.
		 */
		function panic(args: any[]);
	}
	/**
	 * Print panic info and exit.
	 */
	function panic(...args: any[]);
	namespace time {
		export class Time {
			/**
			 * returns unix timestamps (seconds).
			 */
			unix(): number;
		}
	}
	namespace io {
		export class Reader {
			/**
			 * fake method to enable ts-check, or else any value type can be used as io.Reader.
			 */
			read(p: Uint8Array): number;
		}
	}
	namespace fs {
		/**
		 * Mkdir creates a new directory with the specified name and permission bits (before umask).
		 */
		function mkDir(path: string, perm?: number);
		/**
		 * MkdirAll creates a directory named path, along with any necessary parents.
		 * The permission bits perm (before umask) are used for all directories that MkdirAll creates.
		 * If path is already a directory, MkdirAll does nothing.
		 * 
		 * @param path 
		 * @param perm 
		 */
		function mkDirAll(path: string, perm?: number);
		function saveToFile(filePath: string, data: io.Reader): Promise<number>;
		/**
		 * Checks if filePath exists, no matter what kind it is.
		 * @param filePath 
		 * @param types     File types to match.
		 * 
		 * e.g: 'fd' means either file or directory exists. 'fx' means file exists and executable.
		 * it's by default 'fdls'.
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
		function exists(filePath: string, types?: string): boolean;
		/**
		 * Checks if the specified regular file (not directory, not socket, ...) exists.
		 * 
		 * @param filePath 
		 */
		function fileExists(filePath: string): boolean;
		/**
		 * Checks if the specified directory (not regular file, not socket, ...) exists.
		 * 
		 * @param filePath 
		 */
		function dirExists(dirPath: string): boolean;
		/**
		 * Calculates the sha256 sum of a file.
		 * 
		 * @param filePath
		 * 
		 * @returns the digest of the file with each byte represented by two lower-cased letters.
		 */
		function sha256(filePath: string): Promise<string>;
		export class Stat {
			/**
			 * base name of the file.
			 */
			readonly name: string;
			/**
			 * length (in bytes) for regular files.
			 */
			readonly size: number;
			/**
			 * modification time.
			 */
			readonly modTime: time.Time;

			readonly isDir: boolean;
			readonly isRegular: boolean;
		}
		/**
		 * Stat a path to report its info.
		 * @param path 
		 */
		function stat(path: string): Stat;
	}
	namespace archive {
		export class TarReader {
			constructor(input: io.Reader);
			extractTo(dir: string): Promise<void>;
		}
		export class ZipReader {
			constructor(filePath: string);
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
			/**
			 * 200 OK
			 */
			readonly status: string;
			/**
			 * 200
			 */
			readonly statusCode: number;
			/**
			 * statusCode in [200,300).
			 */
			readonly ok: boolean;

			text(): Promise<string>;
			blob(): Promise<ArrayBuffer>;
			reader(): io.Reader;
		}
		function get(url: string): Promise<Response>;
	}
}
export {};
