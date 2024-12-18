#!/bin/bash
set -e

echo "Define variables"
DNS1="8.8.8.8"
DNS2="1.1.1.1"
FREEDNS_URL="https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main"
NGINX_CONF_URL="$FREEDNS_URL/confs/nginx.conf"
FREEDNS_CONF_URL="$FREEDNS_URL/confs/freeDNS.conf"
FREEDNS_SERVICE_URL="$FREEDNS_URL/confs/freeDNS.service"


echo "Stop and disable systemd-resolved"
systemctl stop systemd-resolved
systemctl disable systemd-resolved
pkill freeDNS && echo "freeDNS terminated" || echo "Failed to terminate freeDNS"
pkill nginx && echo "nginx terminated" || echo "Failed to terminate nginx"

echo "Add nameservers if not present"
for DNS in $DNS1 $DNS2; do
    if ! grep -q "nameserver $DNS" /etc/resolv.conf; then
        echo "Adding nameserver $DNS to /etc/resolv.conf"
        echo "nameserver $DNS" | tee -a /etc/resolv.conf > /dev/null
    fi
done

echo "Update package list and install required packages"
apt update
apt install -y ufw nginx wget procps
apt install -y libnginx-mod-stream

echo "Configure UFW"
ufw allow ssh
ufw allow http
ufw allow https
ufw allow 53/udp  
yes | sudo ufw enable

echo "Create directory for freeDNS"
mkdir -p /etc/freeDNS
mkdir -p /usr/local/freeDNS/


echo "Download freeDNS"
wget --no-cache -O /usr/bin/freeDNS "$FREEDNS_URL/scripts/freeDNS.sh"
chmod 755 /usr/bin/freeDNS
wget --no-cache -O /usr/local/freeDNS/freeDNS "$FREEDNS_URL/build/freeDNS"
chmod 755 /usr/local/freeDNS/freeDNS

wget --no-cache -O /etc/systemd/system/freeDNS.service "$FREEDNS_SERVICE_URL"

echo "Backup and replace Nginx configuration"
systemctl stop nginx
mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup
wget --no-cache -O /etc/nginx/nginx.conf "$NGINX_CONF_URL"

echo "Start and enable Nginx"
systemctl start nginx
systemctl enable nginx


echo "Reload systemd and enable freeDNS service"
systemctl daemon-reload
systemctl enable freeDNS
systemctl start freeDNS

