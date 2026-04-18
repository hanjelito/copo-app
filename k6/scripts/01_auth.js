import http from 'k6/http';
import { check } from 'k6'

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
    vus: 100, // 100 usuarios virtuales
    iterations: 100,
}

export default function () {
    const email = `user_${__VU}@copo.com`;

    const res = http.post(`${BASE_URL}/auth/register`,
        JSON.stringify({
            email: email,
            password: '123456',
            name: `User ${__VU}`,
            role: 'passenger',
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    check(res, { 'register ok': (r) => r.status === 201 });
}