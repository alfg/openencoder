FROM golang:1.14-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Outside GOPATH since we're using modules.
WORKDIR /src

# Required for fetching dependencies.
RUN apk add --update --no-cache ca-certificates git nodejs nodejs-npm

# Fetch dependencies to cache.
COPY go.mod go.sum ./
RUN go mod download

# Copy project source files.
COPY . .

# Build static web project.
RUN cd web && npm install && npm run build

# Build.
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix 'static' -v -o /app .

# Final release image.
FROM alfg/ffmpeg:latest

# Set version from CI build.
ARG BUILD_VERSION=${BUILD_VERSION}
ENV VERSION=$BUILD_VERSION

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the project web & executable.
COPY --from=builder /src/web/dist /web/dist
COPY --from=builder /app /app
COPY --from=builder /src/config/default.yml /config/default.yml

EXPOSE 8080

# Perform any further action as an unpriviledged user.
USER nobody:nobody

# Run binary.
ENTRYPOINT ["/app"]
