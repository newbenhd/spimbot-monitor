FROM alpine:latest

# Install required packages
RUN apk add --no-cache bash curl coreutils

# Set working directory
WORKDIR /app

# Copy the binary into the container
COPY build/spimbot-monitor-linux-amd64 /app/spimbot-monitor


# Make sure it’s executable
RUN chmod +x /app/spimbot-monitor

# Default command: run watch on the binary
CMD ["watch", "-n", "3", "/app/spimbot-monitor"]
