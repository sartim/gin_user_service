# Production observability

The service exposes Prometheus metrics at `GET /metrics`. Set `METRICS_TOKEN`
in the deployment secret manager so the endpoint requires a bearer token.
Never commit the token or place it in a container image.

## Prometheus scrape configuration

Store the token in a file readable only by the Prometheus process and reference
that file from the scrape job:

```yaml
scrape_configs:
  - job_name: gin-user-service
    metrics_path: /metrics
    static_configs:
      - targets: [user-service:8000]
    bearer_token_file: /etc/prometheus/secrets/gin-user-service-metrics-token
```

Use the service DNS name or internal load-balancer address as the target. Keep
the metrics route off the public internet; if it must cross an ingress, use
TLS and an explicit allowlist for the Prometheus source network.

## Docker Compose

For a local authenticated scrape, export `METRICS_TOKEN` before starting the
stack. The Compose file passes it through to the application container:

```sh
export METRICS_TOKEN="replace-with-a-local-only-token"
docker compose up -d
curl -H "Authorization: Bearer ${METRICS_TOKEN}" http://localhost:8000/metrics
```

Production deployments should map `METRICS_TOKEN` from a secret-manager
reference rather than an inline environment value.
