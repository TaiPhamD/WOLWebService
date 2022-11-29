#!/bin/bash
sudo systemctl stop wolservice

# check if /opt/wolservice exists 
if [ ! -d "/opt/wolservice" ]; then
    # if not then create the directory
    sudo mkdir -p /opt/wolservice
else
    # remove old binary
    sudo rm /opt/wolservice/wolwebservice
fi
# Copy binary to /opt/wolservice
sudo cp build/dist/wolwebservice /opt/wolservice/
sudo cp wolservice.service /etc/systemd/system/
sudo systemctl enable wolservice
sudo systemctl start wolservice