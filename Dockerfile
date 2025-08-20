# Use Ubuntu as the base image
FROM ubuntu:22.04

# Set maintainer information (optional)
LABEL maintainer="Your Name <your.email@example.com>"

# Update package lists and install basic utilities
RUN apt-get update && \
    apt-get install -y \
    curl \
    wget \
    vim \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && echo "hello world"

# Set working directory
WORKDIR /app

# Copy application files (assuming you have files to copy)
COPY . .

# Expose port (if your application needs it)
EXPOSE 8080

# Set environment variables (if needed)
ENV APP_ENV=production

# Command to run when container starts
CMD ["bash"]