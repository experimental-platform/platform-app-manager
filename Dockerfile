FROM experimentalplatform/ubuntu:latest

COPY platform-app-manager /app-manager

CMD ["dumb-init", "/app-manager", "--port", "80"]

EXPOSE 80
