# docker-bake.hcl

# Build all 3 services by default
group "default" {
  targets = [
    "avrora-app",
    "auth-service",
    "fileserver-service"
  ]
}

target "base" {
  context = "."
  dockerfile = "Dockerfile"
  # Shared args (optional)
  # args = { GO_VERSION = "1.25" }
}

target "avrora-app" {
  inherits = ["base"]
  target = "avrora-app"
  tags = ["avrora-app:latest"]
}

target "auth-service" {
  inherits = ["base"]
  target = "auth-service"
  tags = ["auth-service:latest"]
}

target "fileserver-service" {
  inherits = ["base"]
  target = "fileserver-service"
  tags = ["fileserver-service:latest"]
}