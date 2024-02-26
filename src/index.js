import express from 'express';
import helmet from 'helmet';
import cors from 'cors'
import compression from 'compression';

const PORT = process.env.PORT || 8080;

const app = express();

app.enable('trust proxy');
app.use(cors());
app.use(helmet());
app.use(compression());

app.get('/', (req, res) => {
  const ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
  console.log(ip);
  return res.send(ip + '\n');
});

app.get('/healthz', (req, res) => res.json({ message: "ok"}));
app.use((req, res, next) => res.status(404).json({ message: 'not found'}));
app.use((err, req, res, next) => res.status(500).json({ message: 'error'}));

app.listen(PORT, () => console.log(`Server was started on port http://localhost:${PORT}`));
