# WOLWebService

This is a webservice to send out a Wake on Lan (WOL) command to turn on your computer.

# Installation

1. Run "go build"  in this repository to create WOLWEbService executable
1. Create config.txt based on the config_example.txt and save it in the same path as the executable in 1.

```   
   config.txt content:
    Mypassword1 <---your password
    9998  <--- your server port
    XX:XX:XX:XX:XX:XX <--- MAC address of the PC that you want to turn on (PC has to be setup with Wake on Lan feature enabled)

```
1. You can host this app as a service on raspberry Pi for example:
```
sudo nano /lib/systemd/system/gowol.service
```

Then paste in:

```
[Unit]
Description=GO WOL to turn on PC

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=/home/pi/Documents/WOLWebService/WOLWebService

[Install]
WantedBy=multi-user.target
```
1. Next you can start the service:

```
sudo systemctl enable gowol.service
service gowol start
```

