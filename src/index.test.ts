import request from 'supertest';
import { it, expect, vi } from 'vitest';
import { app as server } from './index.js';
import geoIpLite from 'geoip-lite';

vi.mock('geoip-lite');

const app = request(server);

it('should return with verbose and html', async () => {
	const ip = '127.0.0.1';

	const mockGeoData = {
		range: [123456789, 123456789],
		country: 'US',
		region: 'CO',
		eu: '1',
		timezone: 'America/Chicago',
		city: 'Main',
		ll: [69.2023, -420.846],
		metro: 1,
		area: 0,
	};

	geoIpLite.lookup.mockReturnValueOnce(mockGeoData);

	const response = await app.get(`/?verbose=true`).set('x-forwarded-for', ip);

	expect(response.headers['content-type']).toBe('text/html; charset=utf-8');
	expect(response.text).includes(ip);
	expect(response.text).includes(`<strong>`);
	expect(response.text).includes(mockGeoData.area);
	expect(response.text).includes(mockGeoData.ll);
	expect(geoIpLite.lookup).toHaveBeenCalledTimes(1);
});

it('should return with verbose and json', async () => {
	const ip = '127.0.0.1';

	const mockGeoData = {
		range: [123456789, 123456789],
		country: 'US',
		region: 'CO',
		eu: '1',
		timezone: 'America/Chicago',
		city: 'Main',
		ll: [69.2023, -420.846],
		metro: 1,
		area: 0,
	};

	geoIpLite.lookup.mockReturnValueOnce(mockGeoData);

	const response = await app.get(`/?verbose=true&json=true`).set('x-forwarded-for', ip);

	expect(response.headers['content-type']).toBe('application/json; charset=utf-8');
	expect(response.text).includes(ip);
	expect(response.text).includes(mockGeoData.area);
	expect(response.text).includes(mockGeoData.ll);
	expect(geoIpLite.lookup).toHaveBeenCalledTimes(1);
});

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
