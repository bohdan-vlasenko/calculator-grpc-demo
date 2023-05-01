FROM alpine:latest

# Copy the binary file to the container
COPY calculator_server /

# Run the binary
ENTRYPOINT ["./calculator_server"]