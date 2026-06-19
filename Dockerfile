FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/k3s-sample ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/k3s-sample /k3s-sample

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/k3s-sample"]
