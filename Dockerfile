FROM xena/go:1.12.6 AS build
ENV GOPROXY https://cache.greedo.xeserv.us
COPY . /site
WORKDIR /site
RUN CGO_ENABLED=0 go test -v ./...
RUN CGO_ENABLED=0 GOBIN=/root go install -v ./cmd/site

FROM xena/alpine
EXPOSE 29384
RUN apk add --no-cache bash
WORKDIR /site
COPY --from=build /root/site .
COPY ./cmd/site .
COPY ./static /site/static
COPY ./templates /site/templates
COPY ./blog /site/blog
COPY ./css /site/css
COPY ./app /app
COPY ./app.json .
CMD ./site
