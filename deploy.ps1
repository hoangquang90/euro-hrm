# Thư mục nguồn và đích
$srcDir   = Get-Location
$deployDir = "F:\EuroStack\EuroConn"
$buildOut  = Join-Path $deployDir "euroconn"   # tên binary Linux không có .exe

# Tạo thư mục đích nếu chưa có
if (-Not (Test-Path $deployDir)) {
    New-Item -ItemType Directory -Path $deployDir | Out-Null
}

# Xoá nội dung cũ (nếu cần làm sạch thư mục)
Remove-Item "$deployDir\*" -Recurse -Force -ErrorAction SilentlyContinue

# Set biến môi trường để build cho Linux (cross-compile trên Windows)
$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"   # tránh yêu cầu toolchain C khi cross-compile

# Build binary trực tiếp vào thư mục deploy
Write-Host "Building linux/amd64 to $buildOut ..."
go build -o $buildOut main.go
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed. stop script."
    exit 1
}

# (Tuỳ chọn) Tạo checksum để verify sau khi upload
# Get-FileHash -Path $buildOut -Algorithm SHA256 | ForEach-Object {
#     $_.Hash | Out-File (Join-Path $deployDir "euroconn.sha256") -Encoding ascii
# }

# Danh sách thư mục cần copy
$folders = @("configs", "docs", "init", "internal", "pkg", "resources", "script")
foreach ($folder in $folders) {
    if (Test-Path "$srcDir\$folder") {
        Copy-Item "$srcDir\$folder" -Destination "$deployDir\$folder" -Recurse -Force
    }
}

# Danh sách file cần copy
$files = @("go.mod", "go.sum", "main.go")
foreach ($file in $files) {
    if (Test-Path "$srcDir\$file") {
        Copy-Item "$srcDir\$file" -Destination "$deployDir\$file" -Force
    }
}

Write-Host "Deploy Success: $deployDir"
Write-Host "Binary: $buildOut (linux/amd64, CGO_DISABLED)"