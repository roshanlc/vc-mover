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
echo "Vc-mo#!/bin/bash

# Get user's home dir path and add it dynamically to the last of vc-mover.service
HOME=$(readlink -f ~)

go install github.com/roshanlc/vc-mover@latest && sudo mv ~/go/bin/vc-mover /usr/local/bin
go build ../ && sudo cp ../vc-mover /usr/local/bin

systemctl list-unit-files vc-mover.service && sudo systemctl disable --now vc-mover.service

# write systemd service file
echo "[Unit]
Description=Systemd service file for vc-mover

[Install]
WantedBy=multi-user.target

[Service]
ExecStart=/bin/bash -c '/usr/local/bin/vc-mover --home=$HOME'" > vc-mover.service

sudo chmod 644 vc-mover.service
sudo cp vc-