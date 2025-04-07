import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend } from 'k6/metrics';

export let options = {
    vus: 2000,
    iterations: 200000
};

const targets = [
    { name: 'FastAPI', url: 'http://fastapi:8000/data', metric: 'response_time_fastapi' },
    { name: 'Go', url: 'http://go:8000/data', metric: 'response_time_go' },
    { name: 'NodeJs', url: 'http://node:8000/data', metric: 'response_time_nodejs' }
];

let trends = {};
for (let target of targets) {
    trends[target.name] = new Trend(target.metric);
}

const payload = JSON.stringify({
    name: "Alice",
    email: "alice@example.com",
    age: 28,
    bio: "Backend engineer who loves Go, Python, and NodeJs.",
    interests: ["coding", "reading", "music"]
});

const params = {
    headers: {
        'Content-Type': 'application/json',
    },
};

let currentTarget = __ENV.TARGET || 'FastAPI';
let target = targets.find(t => t.name === currentTarget);

if (!target) {
    throw new Error(`Invalid target: ${currentTarget}`);
}

export default function () {
    const res = http.post(target.url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    trends[target.name].add(res.timings.duration);
    sleep(0.1);
}
