import path from 'path';
// @ts-ignore
import express, { Request, Response, NextFunction } from 'express';
import helmet from 'helmet';
import cors from 'cors';
import dotenv from 'dotenv';
import compression from 'compression';
import { rateLimit as rl } from 'express-rate-limit';

dotenv.config({ path: path.resolve(path.join(process.cwd(), '.env')) });

const PORT = process.env.PORT || 8081;

const app = express();

app.enable('trust proxy');

app.use(
	rl({
		windowMs: 15 * 60 * 1000, // 15 minutes
		limit: 100, // Limit each IP to 100 requests per `window` (here, per 15 minutes).
		standardHeaders: 'draft-7',
		legacyHeaders: false,
		skip: async function (req, _res) {
			const myIp = (req.headers['x-forwarded-for'] || req.socket.remoteAddress).split(', ')[0];
			const myIpWasConnected = myIp === process.env.MY_IP;
			if (myIpWasConnected) console.log(`my ip was connected: ${myIp}`);
			return myIpWasConnected;
		},
		message: (req: Request, res: Response) => {
			const message = 'Too many requests, please try again later?';

			if (req.query.format === 'json' || req.query.json === 'true') {
				return res.status(429).json({ message });
			}

			if (req.get('Content-Type') === 'application/json') {
				return res.status(429).json({ message });
			}

			return res.status(429).send(message + '\n');
		},
	}),
);

app.use(cors());

app.use(helmet());

app.use(compression());

app.use(express.static(path.resolve(path.join(process.cwd(), 'public')), { maxAge: '24h' }));

app.get('/', async (req: Request, res: Response, next: NextFunction) => {
	try {
		const ip = (req.headers['x-forwarded-for'] || req.socket.remoteAddress).split(', ')[0];

		if (req.query.format === 'json' || req.query.json === 'true') {
			return res.status(200).json({ ip });
		}

		if (req.get('Content-Type') === 'application/json') {
			return res.status(200).json({ ip });
		}

		return res.status(200).send(ip + '\n');
	} catch (error) {
		next(error);
	}
});

app.get('/healthz', (req: Request, res: Response, next: NextFunction) => {
	const message = 'ok';

	if (req.get('Content-Type') === 'application/json') {
		return res.status(200).json({ message });
	}

	return res.status(200).send(message + '\n');
});

app.use((req: Request, res: Response, next: NextFunction) => {
	const message = 'not found';

	if (req.get('Content-Type') === 'application/json') {
		return res.status(404).json({ message });
	}

	return res.status(404).send(message + '\n');
});

app.use((error: Error, req: Request, res: Response, next: NextFunction) => {
	const message = 'error';

	if (req.get('Content-Type') === 'application/json') {
		return res.status(500).json({ message });
	}

	return res.status(500).send(message + '\n');
});

const server = app.listen(PORT, () => {
	console.log(`Server was started on http://localhost:${PORT}`);
});

function gracefulShutdown() {
	console.log('Received kill signal, shutting down gracefully.');
	server.close(() => {
		console.log('HTTP server closed.');
		process.exit(0);
	});
}

process.on('SIGINT', gracefulShutdown);

process.on('SIGTERM', gracefulShutdown);

process.on('unhandledRejection', (reason, promise) => {
	console.error('Unhandled Rejection at: ', promise, ' reason: ', reason);
});

export { app };
