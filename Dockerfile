FROM dockerregistry.protorz.net/ubuntu:latest

COPY app-manager /app-manager

CMD ["/app-manager", "--port", "80"]

EXPOSE 80
