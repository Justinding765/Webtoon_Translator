# Start from the official Go image
FROM golang:latest

# Install Python and pip
RUN apt-get update && apt-get install -y python3 python3-pip python3-venv

# Set the working directory inside the container
WORKDIR /app

# Create a Python virtual environment and activate it
RUN python3 -m venv /app/venv
ENV PATH="/app/venv/bin:$PATH"

# Copy the Go and Python dependencies files
COPY go.mod go.sum requirements.txt ./

# Install Go dependencies
RUN go mod download

# Install Python dependencies inside the virtual environment
RUN pip3 install -r requirements.txt

# Copy the Go and Python source code
COPY . .

# Set the environment variable for Google Cloud credentials
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/credentials.json

# Compile the Go application. This includes both server.go and image_translator.go
RUN go build -o server .

# Expose port 5000 to the outside world
EXPOSE 5000

# Run the server program
CMD ["./server"]
