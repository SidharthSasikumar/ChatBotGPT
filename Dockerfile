# Use the official Golang image as a builder stage
FROM golang:1.18 as builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Build the application
RUN go build -o chatbot .

# Use a small image for the final stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/chatbot .

# Run the application
CMD ["./chatbot"]
