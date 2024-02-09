#!/bin/bash

# if go is not installed then exit from the script
if ! command -v go &> /dev/null
then
    echo "*** go could not be on the system, please install it and run the script again ***"
    exit 1
fi

# install the binary and add it to ~/.profile for auto-start at login
go install github.com/roshanlc/vc-mover@latest && echo "exec ~/go/bin/vc-mover" >> ~/.profile

# start the program right away
source ~/.profile &