# Pakai image Golang resmi
FROM golang:1.24

# Bikin folder kerja di container
WORKDIR /app

# Copy dependency file dulu
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build aplikasi Golang kamu
RUN go build -o main .

# Buka port 8080 (ubah sesuai app kamu)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
