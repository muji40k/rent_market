import { check } from 'k6';
import http from 'k6/http';

export const options = {
    stages: [
        {
            duration: '1m',
            target: 500
        },
        {
            duration: '1m',
            target: 500
        },
        {
            duration: '1m',
            target: 0
        }
    ]
}

export default function () {
    const res = http.get(`${__ENV.HOST}`);

    check(res, {
        'accepted': (r) => r.status === 200,
        'rejected': (r) => r.status === 503,
        'other': (r) => r.status !== 200 && r.status !== 503
    });
}

