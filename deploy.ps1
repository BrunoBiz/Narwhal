# Build for Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"

$deployIP = "192.168.18.137" # Edit this

go build -C ./gameserverManager -o ../Narwhal

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed."
    exit 1
}

# Deploy to server - app
scp .\Narwhal gameserver@[$deployIP]:/home/gameserver/Narwhal 

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - .\Narwhal"
    exit 1
}

# CHMOD app
ssh "gameserver@$deployIP" "chmod +x /home/gameserver/Narwhal;"

if ($LASTEXITCODE -ne 0) {
    Write-Host "SSH Failed - CHMOD."
    exit 1
}

Write-Host "Deployment successful!"