FROM golang:1.18 as build

WORKDIR /src
COPY . .
RUN go build

FROM golang:1.18 AS final
ENV TAGS_ADDRESS="0.0.0.0:5000"
EXPOSE 5000
WORKDIR /app
ARG user=rinkudesu
RUN useradd -M -s /bin/false $user
COPY --from=build --chown=$user:$user /src/rinkudesu-tags .
RUN chmod 100 rinkudesu-tags
USER $user
HEALTHCHECK --interval=20s --start-period=5s --retries=3 CMD curl --fail http://localhost:5000/health || exit 1
ENTRYPOINT ["/app/rinkudesu-tags"]