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
    echo -e "${green}Example:${plain}"
    echo -e "  $0 --adddomain example.com --addip 192.168.1.1"
}

# بررسی پارامترهای ورودی
if [[ $# -eq 0 ]]; then
    usage
    exit 1
fi

# اجرای دستورات با پارامترهای ورودی
command="/usr/local/freeDNS"

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
        *)
            echo -e "${red}Unknown option: $1${plain}"
            usage
            exit 1
            ;;
    esac
done

# اجرای دستور نهایی
echo "Executing: $command"
$command
