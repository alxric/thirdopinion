from golang:1.10 AS builder

ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep


WORKDIR $GOPATH/src/thirdopinion
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

ARG RELEASE
ARG BRANCH
ARG COMMIT
ARG GOVER
ENV version $RELEASE
ENV branch $BRANCH
ENV commit $COMMIT
ENV gover $GOVER
ENV user $USER


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=${version} -X main.branch=${branch} -X main.commit=${commit} -X main.goVersion=${gover} -X main.buildUser=${user}" -a -installsuffix nocgo -o /app cmd/thirdopinion/main.go

FROM alpine
COPY --from=builder /app ./
COPY web/template/ /web/template
COPY web/static/ /web/static
EXPOSE 8080
ENTRYPOINT ["./app"]
