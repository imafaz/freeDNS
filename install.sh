#!/bin/bash

# Update DNS settings
echo "Updating DNS settings..."
if ! grep -q "nameserver 8.8.8.8" /etc/resolv.conf; then
    echo "Adding nameserver 8.8.8.8 to /etc/resolv.conf"
    echo "nameserver 8.8.8.8" | sudo tee -a /etc/resolv.conf > /dev/null
fi

if ! grep -q "nameserver 1.1.1.1" /etc/resolv.conf; then
    echo "Adding nameserver 1.1.1.1 to /etc/resolv.conf"
    echo "nameserver 1.1.1.1" | sudo tee -a /etc/resolv.conf > /dev/null
fi

# Update package lists
echo "Updating package lists..."
sudo apt update

# Install required packages
echo "Installing required packages: ufw, nginx, wget..."
sudo apt install -y ufw nginx wget

# Configure UFW
echo "Configuring UFW..."
sudo ufw allow ssh        # Allow SSH
sudo ufw allow http       # Allow HTTP
sudo ufw allow https      # Allow HTTPS
sudo ufw allow 53/udp     # Allow DNS (UDP)
sudo ufw allow 53/tcp     # Allow DNS (TCP, if needed)
sudo ufw enable           # Enable UFW
echo "UFW configuration completed."

# Create directory for freeDNS
echo "Creating directory for freeDNS..."
sudo mkdir -p /etc/freeDNS

# Download required files
echo "Downloading required files..."
wget -O /usr/bin/freeDNS https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/freeDNS
wget -O /etc/freeDNS/freeDNS.conf https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/freeDNS.conf
wget -O /etc/systemd/system/freedns.service https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/freeDNS.service

# Reload systemd
echo "Reloading systemd..."
sudo systemctl daemon-reload

# Enable and start freeDNS service
echo "Enabling and starting freeDNS service..."
sudo systemctl enable freedns.service
sudo systemctl start freedns.service

# Stop Nginx and back up the configuration file
echo "Stopping Nginx and backing up the configuration file..."
sudo systemctl stop nginx
sudo mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup

# Download new Nginx configuration file
echo "Downloading new Nginx configuration file..."
wget -O /etc/nginx/nginx.conf https://raw.githubusercontent.com/imafaz/freeDNS/refs/heads/main/confs/nginx.conf

# Start Nginx and enable it
echo "Starting Nginx and enabling it..."
sudo systemctl start nginx
sudo systemctl enable nginx

echo "Configuration completed successfully."
