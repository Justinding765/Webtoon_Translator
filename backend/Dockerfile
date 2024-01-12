# Start from the official Go image
FROM golang:latest

# Update the package list and install Python, pip, Python virtual environment tools,
# and the necessary libraries for OpenCV and other graphical operations
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-venv \
    libgl1-mesa-glx \  
    libglib2.0-0      

# Set the working directory inside the container to /usr/src/app/backend
WORKDIR /usr/src/app/backend

# Create a Python virtual environment and activate it
RUN python3 -m venv /usr/src/app/backend/venv
ENV PATH="/usr/src/app/backend/venv/bin:$PATH"

# Copy the Go and Python dependencies files
COPY go.mod go.sum requirements.txt ./

# Install Go dependencies
RUN go mod download

# Install Python dependencies inside the virtual environment
RUN pip3 install -r requirements.txt

# Copy the Go and Python source code
COPY . .

# Set the environment variable for Google Cloud credentials
ENV GOOGLE_APPLICATION_CREDENTIALS=/usr/src/app/backend/config/credentials.json


# Expose port 5000 to the outside world
EXPOSE 5000

# Compile the Go application and run the server program when the container launches
CMD go build -o /usr/src/app/backend/server . && /usr/src/app/backend/server
