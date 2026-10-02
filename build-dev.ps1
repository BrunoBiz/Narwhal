# Used in dev - Also builds and deploys the mock server

# Build for Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"

$deployIP = "192.168.18.137"  # Test-server

go build -C ./gameserverManager -o ../gameserver
go build -C ./mockServer -o ../mockServerTest

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed."
    exit 1
}

# Deploy to server - app
scp .\gameserver root@[$deployIP]:/home/gameserver/gameserver 

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - .\gameserver"
    exit 1
}

# Deploy to server - env
scp .\mgr.env root@[$deployIP]:/home/gameserver/mgr.env

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - \mgr.env"
    exit 1
}

# Deploy to server - mockServer
scp .\mockServerTest root@[$deployIP]:/home/gameserver/mockServerTest

if ($LASTEXITCODE -ne 0) {
    Write-Host "Copy failed - .\mockServerTest"
    exit 1
}

# CHMOD both files
ssh "root@$deployIP" "chmod +x /home/gameserver/gameserver; chmod +x /home/gameserver/mockServerTest"

if ($LASTEXITCODE -ne 0) {
    Write-Host "SSH Failed - CHMOD."
    exit 1
}

Write-Host "Deployment successful!"