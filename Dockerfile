FROM golang:1.13-alpine AS build

WORKDIR /go/src/gitlab.com/pztrn/opensaps
COPY . .

RUN go build

FROM alpine:3.10
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/src/gitlab.com/pztrn/opensaps/opensaps /app/opensaps

EXPOSE 25544
ENTRYPOINT [ "/app/opensaps", "-config", "/app/opensaps.yaml" ]
