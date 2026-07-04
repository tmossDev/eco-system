FROM alpine:3.24

WORKDIR /app

ARG APP_NAME

COPY build/bin/${APP_NAME} /app/app
RUN chmod +x /app/app

EXPOSE 8080

CMD ["./app"]
