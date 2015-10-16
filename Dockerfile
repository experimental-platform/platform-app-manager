FROM scratch

COPY platform-app-manager /app-manager

CMD ["/app-manager", "--port", "80"]

EXPOSE 80
