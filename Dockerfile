FROM golang:1.24

# Install air secara global (pakai GOBIN biar masuk ke PATH)
ENV GOBIN=/go/bin
RUN go install github.com/air-verse/air@latest


# Set working directory
WORKDIR /app

# Copy go.mod dulu agar cache efektif
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file project
COPY . .

# Expose port kalau kamu butuh akses dari luar (misal :8080)
EXPOSE 8080

# Jalankan air
ENTRYPOINT ["/go/bin/air"]

