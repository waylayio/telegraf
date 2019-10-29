# Start from the latest golang base image
FROM golang:latest AS builder

# Add Maintainer Info
LABEL maintainer="support@waylay.io"

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/waylay/telegraf
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./

#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /telegraf .
RUN make

# this fails due to internal
# cd $GOPATH
# mkdir -p {src,bin,pkg}
# mkdir -p src/github.com/awesome-org/
# cd src/github.com/awesome-org/
# git clone git@github.com:awesome-you/tool.git # OR: git clone https://github.com/awesome-you/tool.git
# cd tool/
# go get ./...

FROM scratch
COPY --from=builder /telegraf ./
ENTRYPOINT ["./telegraf"]
