FROM golang:1.22.1-alpine3.19 as build
RUN apk add --no-cache make
WORKDIR /app
COPY . .
RUN make build


FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/bin/server /app/server
EXPOSE 8080
ENTRYPOINT ["/app/server"]