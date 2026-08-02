
echo "------------------------Hello to BYTEBEAM----------------"
echo "if you want to transfer a file press 1 else if you want to receive a file press 2"
read input
if [ $input -eq 1 ]
then
    echo "Enter the file path"
    read filePath
    echo "Enter the receiver ip address"
    read ipAddress
    go run ./client/client $filePath $ipAddress
elif [ $input -eq 2 ]
then
    echo "Enter the sender ip address"
    read ipAddress
    go run ./server/server $ipAddress
else
    echo "Invalid input"
fi
