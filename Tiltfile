cfg = config.parse()
backend = cfg.get("backend", "go")

docker_compose("docker-compose.yml")

if backend == "go":
    docker_build(
        "cars-crud-go",
        context="backend/go",
        dockerfile="backend/go/Dockerfile",
        live_update=[
            sync("backend/go", "/app"),
        ],
    )

    local_resource(
        "go-api",
        serve_cmd="cd backend/go && go run ./cmd/api",
        deps=["backend/go"],
        resource_deps=["postgres", "redis", "kafka", "mongo"],
    )