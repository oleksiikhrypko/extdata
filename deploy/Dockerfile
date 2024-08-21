# syntax=docker/dockerfile:1.2

FROM scratch
ARG PROJECT_NAME
COPY migrations migrations
COPY configs/policy.csv /policy.csv
COPY configs/rbac.conf /rbac.conf
COPY configs/jwks.local.json /jwks.local.json
COPY openapi /openapi
COPY api/openapi/api.yaml /openapi/api.yaml
COPY bin/app /app
ENTRYPOINT ["/app", "start", "-c", "/configs/config.yaml"]
