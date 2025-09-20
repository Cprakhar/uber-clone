# Load 'restart_process' extension for live updates with process restarts
load("ext://restart_process", "docker_build_with_restart")

# Apply configurations
k8s_yaml("deployments/k8s/dev/namespace.yaml")
k8s_yaml("deployments/k8s/dev/network.yaml")
k8s_yaml("deployments/k8s/dev/secrets.yaml")

# Deploy the Apache Kafka and Zookeeper
k8s_yaml("deployments/k8s/dev/zookeeper.yaml")
k8s_resource("zookeeper", port_forwards="2181:2181", labels=["messaging"])
k8s_yaml("deployments/k8s/dev/kafka.yaml")
k8s_resource("kafka", port_forwards=["9092:9092", "29092:29092"],
    resource_deps=["zookeeper"], labels=["messaging"]
)


# Deploy the Payment Service
docker_build_with_restart("uber-clone/payment-service:latest", ".",
    dockerfile="services/payment-service/Dockerfile",
    entrypoint=["./main"],
    only=["./services/payment-service", "./shared", "./go.mod", "./go.sum"],
    live_update=[
        sync("services/payment-service", "/app/services/payment-service"),
        sync("shared", "/app/shared"),
    ]
)

k8s_yaml("deployments/k8s/dev/payment-service.yaml")
k8s_resource("payment-service", port_forwards="9200:9200",
    resource_deps=["kafka"],
    labels=["backend"]
)

# Deploy the API Gateway
docker_build_with_restart("uber-clone/api-gateway:latest", ".",
    dockerfile="services/api-gateway/Dockerfile",
    entrypoint=["./main"],
    only=["./services/api-gateway", "./shared", "./go.mod", "./go.sum"],
    live_update=[
        sync("services/api-gateway", "/app/services/api-gateway"),
        sync("shared", "/app/shared"),
    ]
)

k8s_yaml("deployments/k8s/dev/api-gateway.yaml")
k8s_resource("api-gateway", port_forwards="8080:8080",
    resource_deps=["kafka"],
    labels=["backend"]
)

# Deploy the Trip Service
docker_build_with_restart("uber-clone/trip-service:latest", ".",
    dockerfile="services/trip-service/Dockerfile",
    entrypoint=["./main"],
    only=["./services/trip-service", "./shared", "./go.mod", "./go.sum"],
    live_update=[
        sync("services/trip-service", "/app/services/trip-service"),
        sync("shared", "/app/shared"),
    ]
)

k8s_yaml("deployments/k8s/dev/trip-service.yaml")
k8s_resource("trip-service", port_forwards="9000:9000",
    resource_deps=["kafka"],
    labels=["backend"]
)

# Deploy the Driver Service
docker_build_with_restart("uber-clone/driver-service:latest", ".",
    dockerfile="services/driver-service/Dockerfile",
    entrypoint=["./main"],
    only=["./services/driver-service", "./shared", "./go.mod", "./go.sum"],
    live_update=[
        sync("services/driver-service", "/app/services/driver-service"),
        sync("shared", "/app/shared"),
    ]
)

k8s_yaml("deployments/k8s/dev/driver-service.yaml")
k8s_resource("driver-service", port_forwards="9100:9100",
    resource_deps=["kafka"],
    labels=["backend"]
)

# Deploy the Web Service
docker_build("uber-clone/web-service:latest", ".",
    dockerfile="web/Dockerfile",
    only=["./web"],
    live_update=[
        sync("web", "/app"),
    ]
)

k8s_yaml("deployments/k8s/dev/web-service.yaml")
k8s_resource("web-service", port_forwards="3000:3000", labels=["frontend"],
    objects=["stripe:Secret:uber-clone"]
)