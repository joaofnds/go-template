import http from "k6/http";
import { check } from "k6";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: "constant-arrival-rate",
      rate: 100,
      timeUnit: "1s",
      duration: "10s",
      preAllocatedVUs: 100,
      maxVUs: 200,
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95) < 3"],
  },
};

export function setup() {
  http.post("http://localhost:3000/kv/foo/bar");
}

export default function () {
  const res = http.get("http://localhost:3000/kv/foo");

  check(res, {
    "status is 200": (r) => r.status === 200,
    "body is 'bar'": (r) => r.body === "bar",
  });
}
