declare module 'express-serve-static-core' {
	export interface Request {
		headers: {
			'x-forwarded-for': string | string[];
		};
		socket: {
			remoteAddress: string;
		};
	}
}
