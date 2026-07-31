FROM golang:1.26.4 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .
FROM scratch
WORKDIR /
COPY --from=build /app/app /app
COPY --from=build /app/.env /.env
COPY --from=build /app/ddl.sql /ddl.sql
COPY --from=build /app/dml.sql /dml.sql
ENTRYPOINT ["/app"]