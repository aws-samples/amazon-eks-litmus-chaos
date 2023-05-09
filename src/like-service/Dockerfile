## BUILD STAGE
FROM public.ecr.aws/docker/library/golang:1.19-alpine as builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY . .

RUN go build -o app .

## DISTRIBUTION
FROM scratch

COPY --from=builder /app/config/* .
COPY --from=builder /app/app .

CMD ["./app"]
