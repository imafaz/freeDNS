#!/bin/bash


echo "Updating DNS settings..."
if ! grep -q "nameserver 8.8.8.8" /etc/resolv.conf; then
    echo "Adding nameserver 8.8.8.8 to /etc/resolv.conf"
    echo "nameserver 8.8.8.8" | tee -a /etc/resolv.conf > /dev/null
fi

if ! grep -q "nameserver 1.1.1.1" /etc/resolv.conf; then
    echo "Adding nameserver 1.1.1.1 to /etc/resolv.conf"
    echo "nameserver 1.1.1.1" | tee -a /etc/resolv.conf > /dev/null
fi


echo "Updating package lists..."
apt update

echo "Installing required packages: ufw, nginx, wget..."
apt install -y ufw nginx wget


echo "Configuring UFW..."
ufw allow ssh
ufw allow http
ufw allow https
ufw allow 53/udp  
ufw enable          
echo "UFW configuration completed."

echo "Creating directory for freeDNS..."
mkdir -p /etc/freeDNS


echo "Downloading required files..."
wget -O /usr/bin/freeDNS https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/freeDNS
chmod 777 /usr/bin/freeDNS
wget -O /etc/freeDNS/freeDNS.conf https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/freeDNS.conf
wget -O /etc/systemd/system/freedns.service https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/freeDNS.service


echo "Reloading systemd..."
systemctl daemon-reload

systemctl stop systemd-resolved
systemctl disable systemd-resolved

echo "Enabling and starting freeDNS service..."
systemctl enable freedns.service
systemctl start freedns.service


echo "Stopping Nginx and backing up the configuration file..."
systemctl stop nginx
mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup


echo "Downloading new Nginx configuration file..."
wget -O /etc/nginx/nginx.conf https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/nginx.conf


echo "Starting Nginx and enabling it..."
systemctl start nginx
systemctl enable nginx

echo "Configuration completed successfully."