# k3s sample app

A tiny Go HTTP service for testing CI/CD on a k3s cluster.

## Endpoints

- `GET /` returns basic app metadata.
- `GET /healthz` returns liveness status.
- `GET /readyz` returns readiness status.
- `GET /version` returns the configured app version.
- `GET /metrics` returns simple Prometheus-style metrics.

## Run locally

```sh
go run ./cmd/server
```

Then open:

```sh
curl http://localhost:8080/
curl http://localhost:8080/healthz
curl http://localhost:8080/metrics
```

## Test

```sh
go test ./...
```

## Build image

```sh
docker build -t k3s-sample:latest .
```

For GHCR, build and push:

```sh
docker tag k3s-sample:latest ghcr.io/xrellana/k3s-sample:latest
docker push ghcr.io/xrellana/k3s-sample:latest
```

## Deploy to k3s

```sh
kubectl apply -f k8s/app.yaml
kubectl -n k3s-sample get pods,svc,ingress
```

If you use the included Ingress, point `k3s-sample.localhost` to your k3s node IP.
