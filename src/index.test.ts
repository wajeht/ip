import request from 'supertest';
import { it, expect } from 'vitest';
import { app as server } from './index.js';

const app = request(server);

it('should be able to get an ip address', async () => {
	const response = await app.get('/').set('Content-Type', 'application/json');
	expect(response.status).toBe(200);
	expect(response.headers['content-type']).toBe('application/json; charset=utf-8');
	expect(response.body).toHaveProperty('ip');
});

it('should be able to ping healthz end point', async () => {
	const response = await app.get('/healthz');
	expect(response.status).toBe(200);
	expect(response.text).include('ok');
});

it('should be able to hit not found', async () => {
	const response = await app.get('/not-found');
	expect(response.status).toBe(404);
	expect(response.text).include('Not found');
});
