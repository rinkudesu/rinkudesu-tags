FROM golang:1.18-alpine as build

WORKDIR /src
COPY . .
RUN go build

FROM golang:1.18-alpine AS final
ENV TAGS_ADDRESS="0.0.0.0:5000"
EXPOSE 5000
WORKDIR /app
ARG user=rinkudesu
RUN adduser -D -H -s /bin/false $user
RUN addgroup $user $user
COPY --from=build --chown=$user:$user /src/rinkudesu-tags .
RUN chmod 100 rinkudesu-tags
USER $user
ENTRYPOINT ["/app/rinkudesu-tags"]