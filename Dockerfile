FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -tags timetzdata

FROM scratch

COPY --from=app-builder /go/src/app/clean-dir-on-request /clean-dir-on-request

ENTRYPOINT ["/clean-dir-on-request"]
