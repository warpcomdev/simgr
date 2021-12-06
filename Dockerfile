FROM node:14-alpine as node_builder
WORKDIR /app
COPY package.json package-lock.json /app/
RUN  npm ci
COPY .eslintrc.cjs .prettierrc svelte.config.js tsconfig.json /app/
COPY static /app/static
COPY src /app/src
RUN  npm run build

FROM golang:1.17 AS builder
WORKDIR /app
COPY go.mod go.sum *.go /app/
COPY internal /app/internal
COPY --from=node_builder /app/build /app/build
RUN  CGO_ENABLED=0 go build

FROM busybox:1.33
WORKDIR /app
COPY --from=node_builder /app/build /app/build
COPY --from=builder /app/simgr /app/simgr
ENTRYPOINT ["/app/simgr"]
