#!/bin/bash

# Install the binary
go install github.com/roshanlc/vc-mover@latest

# create the dir if not exists
mkdir -p ~/.config/systemd/user

# Write the unit file directly
echo "
[Unit]
Author=Roshan Lamichhane
Description=Systemd service file for vc-mover

[Install]
WantedBy=default.target

[Service]
ExecStart=/bin/bash -c ~/go/bin/vc-mover
" > ~/.config/systemd/user/vc-mover.service

# Copy the systemd unit file
# cp vc-mover.service ~/.config/systemd/user

# Disable the existing vc-mover service and re-enable it again
systemctl --user list-unit-files vc-mover.service && systemctl --user disable --now vc-mover.service
systemctl --user enable --now vc-mover.service

echo "Vc-mover is now up and running"