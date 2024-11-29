#!/bin/bash

red='\033[0;31m'
green='\033[0;32m'
yellow='\033[0;33m'
plain='\033[0m'

usage() {
    echo -e "${green}Usage: $0 [options]${plain}"
    echo -e "${yellow}Options:${plain}"
    echo -e "  -h, --help        Show this help message"
    echo -e "  -v                Show version"
    echo -e "  -server <IP>     Set DNS server listen IP"
    echo -e "  -port <port>     Set DNS server listen port"
    echo -e "  -adddomain <domain>  Add domain"
    echo -e "  -addip <IP>      Add IP"
    echo -e "  -deldomain <domain>  Delete domain"
    echo -e "  -delip <IP>      Delete IP"
    echo -e "  -start           Start freeDNS service"
    echo -e "  -stop            Stop freeDNS service"
    echo -e "  -status          Check the status of freeDNS service"
    echo -e "  -install         Install freeDNS"
    echo -e "  -uninstall       Uninstall freeDNS"
    echo -e "${green}Example:${plain}"
    echo -e "  $0 --adddomain example.com --addip 192.168.1.1"
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
        [[ $# -eq 0 ]] && start
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
    rm -f /etc/systemd/system/freeDNS.service
    systemctl daemon-reload
    rm -rf /etc/freeDNS/
    rm -rf /usr/local/freeDNS
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
            command="$command -v"
            shift
            ;;
        -server)
            shift
            command="$command --server $1"
            shift
            ;;
        -port)
            shift
            command="$command --port $1"
            shift
            ;;
        -adddomain)
            shift
            command="$command --adddomain $1"
            shift
            ;;
        -addip)
            shift
            command="$command --addip $1"
            shift
            ;;
        -deldomain)
            shift
            command="$command --deldomain $1"
            shift
            ;;
        -delip)
            shift
            command="$command --delip $1"
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
