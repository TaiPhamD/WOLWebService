#!/bin/bash
sudo systemctl stop wolservice
sudo cp unifi_udm_scripts/wolservice.service /etc/systemd/system/
sudo systemctl enable wolservice
sudo systemctl start wolservice