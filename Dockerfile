# Build stage for Go application
FROM golang:1.19-alpine 

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite

# Set working directory
WORKDIR /app

# Copy go mod files
COPY . .

RUN apk add --no-cache python3 py3-pip && \
    pip3 install pandas && \
    python3 convert_to_sqlite.py && \
    rm airports.csv convert_to_sqlite.py && \
    apk del python3 py3-pip

# Download dependencies
RUN go mod tidy

# Build the application
RUN go build -o airport-map .

# Install runtime dependencies
RUN apk add --no-cache sqlite

# Expose port
EXPOSE 8080

# Run the application
CMD ["./airport-map"] 