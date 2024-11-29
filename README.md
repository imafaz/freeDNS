# freeDNS
## freeDNS is a small project aimed at bypassing software restrictions. It’s written in Go and uses Nginx for reverse proxy.

> **Disclaimer:** This project is intended solely for personal learning and communication purposes. Please refrain from using it for any illegal activities or in a production environment.

**If this project is helpful to you, you may wish to give it a**:star2:

### Table of contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [License](#license)
- [Contributors](#contributors)

### Features
- IP restriction for the DNS server
- Adding more domains
- Start, stop, and check the status of the DNS service
- Add and delete domains and IPs dynamically

### Prerequisites
- Ubuntu 20 or higher version

### Installation
```bash
bash <(curl -Ls https://raw.githubusercontent.com/imafaz/freeDNS/main/install.sh)
```

### Usage
To use the freeDNS script, you can run it with various options. Below are the available commands and their descriptions:

```bash
./freeDNS [options]
```

#### Options:
- `-h, --help`        Show this help message
- `-v`                Show version
- `-server <IP>`     Set DNS server listen IP
- `-port <port>`     Set DNS server listen port
- `-adddomain <domain>`  Add a domain to the DNS server
- `-addip <IP>`      Add an IP address to the DNS server
- `-deldomain <domain>`  Delete a domain from the DNS server
- `-delip <IP>`      Delete an IP address from the DNS server
- `-start`           Start the freeDNS service
- `-stop`            Stop the freeDNS service
- `-status`          Check the status of the freeDNS service

#### Example:
To add a domain , you can use the following command:
```bash
./freeDNS --adddomain google.com
```



### License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

### Contributors
Feel free to contribute to the project by submitting issues or pull requests!