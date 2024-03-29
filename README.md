# WOLWebService
## summary
This webservice could operate in two modes:
  - Master: This mode is meant for an always on device like a raspberry pi or a unix based router like the UDM PRO/SE
      - It could relay Wake-On-Lan commands to your PC to wake it up.
      - Query slave system to get current OS information
      - Reboot slave system and change its UEFI BootNext variable with the specified BootID
  - Slave: This mode is meant for the client computer that would fulfill restart requests or OS query

A typical use case would be to install this server app on an always on device like a raspberry pi/Linux router then set it up in Master mode. Then configure config.json to include information about all your LAN devices that you want to be able to send WoL packets, restart command, or query OS info. You do not need to install this software on the WoL targeted system unless you want to support other APIs like suspend or restarting to a different OS.

If you don't have an always on device then you could install this service directly on the computer and you would just lose out the WoL functionality. 


## API

These APIs were built with security in mind where it will make sure your api_key matches the server's config.json before it processes any API call. Your local PC client information is also secured since you are only sending alias information instead of an actual IP/MAC address of the desired client system. There is also an internal rate limiter that make sure you can't spam the APIs.

- POST /api/wol   - This API sends a Wake-On-Lan packet to your client PC defined by the alias param. Only available in Master mode.
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1"
}
```
- POST /api/os   - This API query OS type from your client PC defined by the alias param
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1"
}
```
- POST /api/restart   - This API restarts the client PC based on the OS parameter
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1",
        "os": "ubuntu" <---Optional parameter. If no OS is specified then it will restart without changing the UEFI BootNext variable.
}
```
- POST /api/suspend   - This API suspends the client PC 
```
{
	"api_key": "4235sdfadf",
	"alias": "pc1"
}
```
## Automation

Once you have the server setup you can easily send HTTPS payload to your server using any of the following:
- Google Assistant --> IFTTT --> your webserver
- Siri Shortcuts --> your webserver
    - Note: Siri Shortcuts won't allow you to connect to a self-signed HTTPS so you should get a real HTTPS cert for free via [LetsEncrypt](https://letsencrypt.org)
    
#### Example of Siri shortcut to send WoL payload    
Here is an example of a siri shortcut:

<img width="713" alt="Screenshot 2022-11-28 at 6 46 23 PM" src="https://user-images.githubusercontent.com/10516699/204426683-538dd29b-d032-4128-a9a3-0e8bc9f00de6.png">



## Config file explanation

Master config
```
{
    "master": true,
    "tls": false,
    "port": "9991", <--- Listening port of web appp server
    "rate_limit": 1, <--- limit max api request per sec
    "rate_burst": 1, <--- limit api request max burst rate
    "api_key": "my_secure_password", <--- Password to authorize api calls
    "fullchain": "path_to/certs/fullchain.pem", <--- optional TLS certs for HTTPS hosting
    "priv_key": "path_to/certs/privkey.pem",<--- optional TLS certs for HTTPS hosting
    "clients:": [ <--- Add all PCs on your LAN that you want WOL control here
        {
            "alias": "client1", <--- alias used to select the right PC . aka mapping to IP/MAC info.
            "ip": "192.168.2.23:9991",
            "mac": "00:00:00:00:00:00"
        },
        {
            "alias": "client2",
            "ip": "192.168.0.27", <--- if no port specified then it assumes the same port as the master's config
            "mac": "aa:aa:aa:aa:aa:aa" 
        }
    ]
}
```

The slave config is almost like the master config except the ip/mac information isn't needed. If you are using the restart to specific OS then you will need to model OS information including the boot_id which is obtained via any standadard UEFI boot manager like efibootmgr on linux.


Slave config
```
{
    "master": false,
    "tls": false,
    "port": "9991",
    "rate_limit": 1, <--- limit max api request per sec
    "rate_burst": 1, <--- limit api request max burst rate
    "api_key": "my_secret_key",
    "clients:": [
        {
            "alias": "client1",
            "os": [
                {
                    "name": "Windows",
                    "boot_id": "0000" <-- UEFI boot id could be obtained from efibootmgr (linux app)
                },
                {
                    "name": "ubuntu",
                    "boot_id": "0002"
                }
            ]
        },
        {
            "alias": "client2",
            "os": [
                {
                    "name": "Windows",
                    "boot_id": "0001"
                }
            ]
        }
    ]
}
```

# Build from source
## Linux
### Prerequisite
- Linux distro with systemd support. If you don't have it then you just have to modify the installation script yourself and edit 1 line for the suspend code [here](https://github.com/TaiPhamD/WOLWebService/blob/8137ca66b9ac6d4dea3cd1b5e4d359f3b6c33a92/server/util/util_linux.go#L12) to not rely on systemctl.
- Go lang compiler 
### Build step
- ./build_linux.sh
   - tip if you need to cross compile set env GOOS=linux GOARCH=arm64 if you want to compile for an arm device like a raspberry pi
### Install step
- ./install_linux.sh (It will install a systemd service named wolservice and start it but it wont work yet until you setup a config.json)
- setup config.json based on the examples from [master config](https://github.com/TaiPhamD/WOLWebService/blob/master/config_sample_master.json)
or [slave config](https://github.com/TaiPhamD/WOLWebService/blob/master/config_sample_slave.json)
- copy config.json to /opt/wolservice/config.json
- sudo systemctl restart wolservice
## Windows
### Prerequisite
- MSYS2 for GCC tool chain (https://www.msys2.org/). Install the following packages:
    - pacman -S --needed base-devel mingw-w64-x86_64-toolchain
    - pacman -S cmake msys2-w32api-headers msys2-w32api-runtime
- Golang 1.16+
- CMake v3+
### Build step
- ./build_windows.bat
### Install step
- ./install_windows.bat (It will install a Windows Service called WOLServerService
- setup config.json based on the examples from [master config](https://github.com/TaiPhamD/WOLWebService/blob/master/config_sample_master.json)
or [slave config](https://github.com/TaiPhamD/WOLWebService/blob/master/config_sample_slave.json)
- copy config.json to C:\wolservice\config.json
- Restart the WOLServerService from services.msc
