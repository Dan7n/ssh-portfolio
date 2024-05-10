# syntax=docker/dockerfile:1

FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /root


RUN mkdir /var/run/sshd
CMD ["/usr/sbin/sshd", "-D"]

# Download Go modules
COPY . .
RUN go mod download

# Copy source code
COPY *.go ./

# Build the image
RUN CGO_ENABLED=0 GOOS=linux go build -C /root/cmd/api -o /root/ssh-portfoloio

# Run the container
FROM alpine:latest AS runner
WORKDIR /root
COPY --from=builder /root/ssh-portfoloio .
EXPOSE 2222

# Start the application
CMD ["/root/ssh-portfoloio"]

