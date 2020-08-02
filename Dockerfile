FROM registry.access.redhat.com/ubi8/go-toolset as builder

USER root

# Setup go environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

# Change working directory
WORKDIR $GOPATH/src/verifypass-go/

# Install dependencies
ENV GO111MODULE=on
COPY . ./
RUN if test -e "go.mod"; then go build ./...; fi

RUN go build -o $GOPATH/bin/verifypass github.com/vivshankar/verifypass-go/cmd/verifypass

FROM registry.access.redhat.com/ubi8/ubi-minimal

COPY --from=builder /go/bin/verifypass /usr/local/bin/
COPY --from=builder /go/src/verifypass-go/public ./public/

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

CMD ["verifypass"]