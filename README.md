# freeDNS
## freeDNS is a small project aimed at bypassing software restrictions. It’s written in Go and uses Nginx for reverse proxy.

> **Disclaimer:** This project is intended solely for personal learning and communication purposes. Please refrain from using it for any illegal activities or in a production environment.


**If this project is helpful to you, you may wish to give it a**:star2: **to support future updates and feature additions!**


### Notes:
- This project is newly released and may contain bugs; it is not recommended for organizational use.
- Supporting this project with a star will help in future updates and feature additions.
- The project currently works on specific operating systems and AMD architecture, but support for more operating systems and architectures will be added soon.


### Table of contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Todo](#todo)
- [License](#license)
- [Contributors](#contributors)

### Features
- IP restriction for the DNS server
- Adding more domains
- Start, stop, and check the status of the DNS service
- Add and delete domains and IPs dynamically
- Enable specific domains and IP restrictions
- List all domains and allowed IPs
- Install, uninstall, and update the service

### Prerequisites
- Ubuntu 20 or higher version

### Installation
```bash
bash <(curl -Ls https://raw.githubusercontent.com/imafaz/freeDNS/main/scripts/install.sh)
```

### Usage
To use the freeDNS script, you can run it with various options. Below are the available commands and their descriptions:

```bash
freeDNS [options]
```

#### Options:
- `-h, --help`                    Show this help message
- `-v`                            Show version
- `-dns-server-ip <IP>`          Set DNS server listen IP
- `-dns-server-port <port>`      Set DNS server listen port
- `-add-domain <domain>`         Add a domain to the DNS server
- `-add-ip <IP>`                 Add an IP address to the DNS server
- `-delete-domain <domain>`      Delete a domain from the DNS server
- `-delete-ip <IP>`              Delete an IP address from the DNS server
- `-start-server`                 Start the DNS server
- `-reverse-proxy-ip <IP>`       Set the reverse proxy Nginx IP
- `-enable-specific-domains <yes/no>` Enable specific domains
- `-enable-ip-restrictions <yes/no>` Enable IP restrictions
- `-list-domains`                 Show all domains
- `-list-ips`                     Show all allowed IPs
- `-list-configs`                 Show all configs
- `-start`                        Start the freeDNS service
- `-stop`                         Stop the freeDNS service
- `-status`                       Check the status of the freeDNS service
- `-install`                      Install freeDNS
- `-uninstall`                    Uninstall freeDNS
- `-update`                       Update freeDNS

#### Example:
To add a domain, you can use the following command:
```bash
freeDNS -add-domain example.com
```

### TODO:
- [ ] Add restricted domains from GitHub
- [ ] Add wildcard domain support (including all subdomains)
- [ ] Write a UI panel for this project

### License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

### Contributors
Feel free to contribute to the project by submitting issues or pull requests!