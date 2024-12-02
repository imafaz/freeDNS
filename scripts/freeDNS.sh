#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

usage() {
    echo -e "${green}Usage: freeDNS [options]${plain}"
    echo -e "${yellow}Options:${plain}"
    echo -e "  -h, --help                    Show this help message"
    echo -e "  -v                            Show version"
    echo -e "  -dns-server-ip <IP>          Set DNS server listen IP"
    echo -e "  -dns-server-port <port>      Set DNS server listen port"
    echo -e "  -add-domain <domain>         Add domain"
    echo -e "  -add-ip <IP>                 Add IP"
    echo -e "  -delete-domain <domain>      Delete domain"
    echo -e "  -delete-ip <IP>              Delete IP"
    echo -e "  -start-server                 Start DNS server"
    echo -e "  -reverse-proxy-ip <IP>       Reverse proxy nginx IP"
    echo -e "  -enable-specific-domains <yes/no> Enable specific domains"
    echo -e "  -enable-ip-restrictions <yes/no> Enable IP restrictions"
    echo -e "  -list-domains                 Show all domains"
    echo -e "  -list-ips                     Show all allowed IPs"
    echo -e "  -list-configs                 Show all configs"
    echo -e "  -start           Start freeDNS service"
    echo -e "  -stop            Stop freeDNS service"
    echo -e "  -status          Check the status of freeDNS service"
    echo -e "  -install         Install freeDNS"
    echo -e "  -uninstall       Uninstall freeDNS"
    echo -e "${green}Example:${plain}"
    echo -e "  freeDNS --add-domain example.com"
}

confirm() {
    if [[ $# -gt 1 ]]; then
        echo && read -p "$1 [Default: $2]: " temp
        if [[ -z "${temp}" ]]; then
            temp=$2
        fi
    else
        read -p "$1 [y/n]: " temp
    fi
    if [[ "${temp,,}" == "y" ]]; then
        return 0
    else
        return 1
    fi
}

install() {
    echo "Installing freeDNS..."
    bash <(curl -Ls https://raw.githubusercontent.com/imafaz/freeDNS/main/scripts/install.sh)
    
    if [[ $? -eq 0 ]]; then
        echo "Installation successful."
    else
        echo -e "${red}Installation failed.${plain}"
    fi
}

uninstall() {
    confirm "Are you sure you want to uninstall freeDNS?" "n"
    if [[ $? -ne 0 ]]; then
        return 0
    fi
    echo "Uninstalling freeDNS..."
    systemctl stop freeDNS
    systemctl disable freeDNS
    systemctl stop nginx
    systemctl disable nginx
    apt remove nginx -y
    rm -f /etc/systemd/system/freeDNS.service
    systemctl daemon-reload
    rm -rf /etc/freeDNS/
    rm -rf /usr/local/freeDNS
    rm -f /var/log/freeDNS.log
    echo "Uninstallation complete."
}

if [[ $# -eq 0 ]]; then
    usage
    exit 1
fi

command="/usr/local/freeDNS/freeDNS"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -v)
            command="$command --version"
            shift
            ;;
        -dns-server-ip)
            shift
            command="$command --dns-server-ip $1"
            shift
            ;;
        -dns-server-port)
            shift
            command="$command --dns-server-port $1"
            shift
            ;;
        -add-domain)
            shift
            command="$command --add-domain $1"
            shift
            ;;
        -add-ip)
            shift
            command="$command --add-ip $1"
            shift
            ;;
        -delete-domain)
            shift
            command="$command --delete-domain $1"
            shift
            ;;
        -delete-ip)
            shift
            command="$command --delete-ip $1"
            shift
            ;;
        -start-server)
            command="$command --start-server"
            shift
            ;;
        -reverse-proxy-ip)
            shift
            command="$command --reverse-proxy-ip $1"
            shift
            ;;
        -enable-specific-domains)
            shift
            command="$command --enable-specific-domains $1"
            shift
            ;;
        -enable-ip-restrictions)
            shift
            command="$command --enable-ip-restrictions $1"
            shift
            ;;
        -list-domains)
            command="$command --list-domains"
            shift
            ;;
        -list-ips)
            command="$command --list-ips"
            shift
            ;;
        -list-configs)
            command="$command --list-configs"
            shift
            ;;
        -start)
            echo "Starting freeDNS service..."
            systemctl start freeDNS
            exit 0
            ;;
        -stop)
            echo "Stopping freeDNS service..."
            systemctl stop freeDNS
            exit 0
            ;;
        -status)
            if systemctl is-active --quiet freeDNS; then
                echo -e "${green}freeDNS is running.${plain}"
            else
                echo -e "${red}freeDNS is not running.${plain}"
            fi
            exit 0
            ;;
        -install)
            install
            exit 0
            ;;
        -uninstall)
            uninstall
            exit 0
            ;;
        *)
            echo -e "${red}Unknown option: $1${plain}"
            usage
            exit 1
            ;;
    esac
done

echo "Executing: $command"
$command
