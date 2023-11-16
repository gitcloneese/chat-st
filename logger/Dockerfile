FROM golang:1.20 AS builder

WORKDIR /src

COPY . .
 
WORKDIR /src/logger

RUN git config --global url."https://oauth2:glpat-iehtnh3GhsW-4Gp3eAyo@gitlab.com".insteadOf "https://gitlab.com" \
    && go env -w GOPRIVATE=gitlab.com \
    && go env -w GOPROXY=https://goproxy.cn,direct && go env -w GOSUMDB="sum.golang.google.cn" \
    && go mod edit -replace xy3-proto=gitlab.com/firerocksg/xy3-proto/xy3-proto@develop
# To workaround GOPROXY unexpected EOF we run it a second time to repull failed packages
RUN --mount=type=cache,target=/root/.cache/go-build go mod tidy || go mod tidy 
RUN --mount=type=cache,target=/root/.cache/go-build make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y


RUN mkdir -p /opt/service/bin
WORKDIR /opt/service

COPY --from=builder /src/logger/bin  /opt/service/bin/

ENTRYPOINT ["/opt/service/bin/cmd", "-conf", "/opt/service/configs/base,/opt/service/configs/args"]
