#!/bin/bash

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
sudo cp vc-mover.service /etc/systemd/system && sudo systemctl enable --now vc-mover.service