FROM code.pztrn.name/containers/mirror/golang:1.18.3-alpine AS build

WORKDIR /go/src/gitlab.com/pztrn/opensaps
COPY . .

RUN go build

FROM code.pztrn.name/containers/mirror/alpine:3.16.0
LABEL maintainer "Stanislav N. <pztrn@pztrn.name>"

COPY --from=build /go/src/gitlab.com/pztrn/opensaps/opensaps /app/opensaps

EXPOSE 25544
ENTRYPOINT [ "/app/opensaps", "-config", "/app/opensaps.yaml" ]
