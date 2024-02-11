#!/bin/bash

# if go is not installed then exit from the script
if ! command -v go &> /dev/null
then
    echo "*** go could not be on the system, please install it and run the script again ***"
    exit 1
fi

# install the binary and add it to ~/.profile for auto-start at login
go install github.com/roshanlc/vc-mover@latest

# append to ~/.profile only in the case of first installation
if ! grep -q 'exec ~/go/bin/vc-mover' ~/.profile
then
    echo "exec ~/go/bin/vc-mover" >> ~/.profile
fi

# start the program right away
exec ~/go/bin/vc-mover &
echo "Vc-mover has been started in background"