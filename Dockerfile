FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/ ./cmd/...

FROM alpine:3.20
COPY --from=build /out/api /api
COPY --from=build /out/migrate /migrate
EXPOSE 8080
ENTRYPOINT ["/api"]
