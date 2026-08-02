Write-Host "------------------------Hello to BYTEBEAM----------------"
Write-Host "If you want to transfer a file press 1, else if you want to receive a file press 2"

$inputChoice = Read-Host "Enter your choice"

if ($inputChoice -eq "1") {
    Write-Host "Beginning file transfer..."
    cd client
    ./client.exe $filePath $ipAddress
}
elseif ($inputChoice -eq "2") {
    Write-Host "Starting file receiver..."
    cd server
    ./server.exe 
}
else {
    Write-Host "Invalid input"
}