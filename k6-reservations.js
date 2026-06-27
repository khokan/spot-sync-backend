import http from 'k6/http';
import { check, sleep } from 'k6';

http.setResponseCallback(http.expectedStatuses(201, 409));

const BASE_URL = __ENV.BASE_URL || 'http://127.0.0.1:8080';
const TOKEN = __ENV.TOKEN || '';
const ZONE_ID = Number(__ENV.ZONE_ID || '4');

if (!TOKEN) {
  throw new Error('TOKEN environment variable is required');
}

export const options = {
  scenarios: {
    reservation_create: {
      executor: 'constant-arrival-rate',
      rate: Number(__ENV.RATE || '40'),
      timeUnit: '1s',
      duration: __ENV.DURATION || '30s',
      preAllocatedVUs: Number(__ENV.PRE_VUS || '100'),
      maxVUs: Number(__ENV.MAX_VUS || '300'),
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'],
    http_req_duration: ['p(95)<800'],
  },
};

function uniquePlate() {
  // Keep <= 15 chars to satisfy backend validation max=15.
  const vu = (__VU % 100).toString().padStart(2, '0');
  const iter = (__ITER % 10000).toString().padStart(4, '0');
  const ts = Date.now().toString().slice(-5);
  return `KK${vu}${iter}${ts}`;
}

export default function () {
  const url = `${BASE_URL}/api/v1/reservations`;
  const payload = JSON.stringify({
    zone_id: ZONE_ID,
    license_plate: uniquePlate(),
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${TOKEN}`,
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 201 or 409': (r) => r.status === 201 || r.status === 409,
  });

  // Small pause to reduce unrealistically bursty client behavior.
  sleep(0.05);
}
