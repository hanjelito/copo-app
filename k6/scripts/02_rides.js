import http from 'k6/http';
import { check } from 'k6'

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
    vus: 10,
    iterations: 100,
}

export default function () {
    const email = `user_${__VU}@copo.com`;

    const res = http.post(`${BASE_URL}/auth/login`,
        JSON.stringify({
            email: email,
            password: '123456',
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    const token = res.json('access_token');

    const ride = http.post(`${BASE_URL}/rides`,
        JSON.stringify({
            origin: 'Madrid',
            destination: 'Barcelona',
            departure: '2026-06-01T08:00:00Z',
            seats: 3
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        }
    )

    check(ride, { 'ride created': (r) => r.status === 201 })
}