cfg = config.parse()
backend = cfg.get("backend", "go")

docker_compose("docker-compose.yml")

if backend == "go":
    local_resource(
        "go-api",
        serve_cmd="cd backend/go && go run ./cmd/api",
        deps=["backend/go"],
        resource_deps=["postgres", "redis", "kafka", "mongo"],
    )

local_resource(
    "react-dev",
    serve_cmd="cd frontend/react && npm run dev -- --port 3000 --host",
    deps=["frontend/react/src"],
    resource_deps=["go-api"],
    links=["http://localhost:3000"],
)