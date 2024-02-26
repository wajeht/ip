import express from 'express';

const app = express();
const PORT = process.env.PORT || 8080;

app.enable('trust proxy');

app.get('*', (req, res) => {
  const ip = req.headers['x-forwarded-for'] || req.socket.remoteAddress;
  return res.send(ip + '\n');
});

app.listen(PORT, () => console.log(`Server running on port http://localhost:${PORT}`));
