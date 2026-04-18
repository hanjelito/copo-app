import http from 'k6/http'
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
            seats: 3,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        }
    );

    const rideID = ride.json('id');

    const booking = http.post(`${BASE_URL}/bookings`,
        JSON.stringify({
            ride_id: rideID,
        }),
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            }
        }
    );

    const bookingID = booking.json('id')

    const bookingList = http.get(`${BASE_URL}/bookings/me`,
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            }
        }
    );

    check(bookingList, { 'bookings listed': (r) => r.status === 200 });

    const bookingDelete = http.del(`${BASE_URL}/bookings/${bookingID}`, null,
        {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            }
        }
    )

    check(bookingDelete, { 'booking cancelate': (r) => r.status === 204 })

}
