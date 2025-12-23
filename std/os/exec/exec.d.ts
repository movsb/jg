declare global {
	namespace exec {
		class Command {
			constructor(cmd: string);
		}
		const command: {
			new (cmd: string): Command;
		};
	}
}
export {};
