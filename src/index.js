import path from 'path';
import express from 'express';
import helmet from 'helmet';
import cors from 'cors';
import compression from 'compression';
import { rateLimit as rl } from 'express-rate-limit';

const PORT = process.env.PORT || 8080;

const app = express();

app.enable('trust proxy');

app.use(
	rl({
		windowMs: 15 * 60 * 1000, // 15 minutes
		limit: 100, // Limit each IP to 100 requests per `window` (here, per 15 minutes).
		standardHeaders: 'draft-7',
		legacyHeaders: false,
		message: (req, res) => {
			return res.json({
				message: 'Too many requests, please try again later?',
			});
		},
	}),
);

app.use(cors());

app.use(helmet());

app.use(compression());

app.use(express.static(path.resolve(path.join(process.cwd(), 'public')), { maxAge: '24h' }));

app.get('/', (req, res) => {
	const ip = (req.headers['x-forwarded-for'] || req.socket.remoteAddress).split(', ')[0];
	if (req.get('Content-Type') === 'application/json') {
		return res.json({ ip });
	}
	return res.send(ip + '\n');
});

app.get('/healthz', (req, res) => res.json({ message: 'ok' }));

app.use((req, res, _next) => res.status(404).json({ message: 'not found' }));

app.use((err, req, res, _next) => res.status(500).json({ message: 'error' }));

app.listen(PORT, () => console.log(`Server was started on port http://localhost:${PORT}`));

export { app };
