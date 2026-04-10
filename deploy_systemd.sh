#!/bin/bash
# Deploy hoyj_fs API as a systemd service
# Usage: ./deploy_systemd.sh

set -e

SERVER="r4wm@45.79.111.243"

echo "==> Connecting to server and building Go binary..."

ssh "$SERVER" bash <<'ENDSSH'
set -e

export PATH="/usr/local/go/bin:$PATH"
PROJECT_DIR="/home/r4wm/github/hoyj_fs"
BINARY_NAME="hoyj-api"

cd "$PROJECT_DIR"

echo "==> Pulling latest code..."
git pull || echo "Git pull failed or not a git repo, continuing..."

echo "==> Building Go binary..."
cd api
go build -o "$BINARY_NAME" server.go
mv "$BINARY_NAME" /home/r4wm/

echo "==> Binary built at /home/r4wm/$BINARY_NAME"

echo "==> Creating service file in home directory..."
cat > /home/r4wm/hoyj-api.service <<'SERVICEEOF'
[Unit]
Description=HOYJ File Server API
After=network.target redis.service
Wants=redis.service

[Service]
Type=simple
User=r4wm
Group=r4wm
WorkingDirectory=/home/r4wm
ExecStart=/home/r4wm/hoyj-api
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin

[Install]
WantedBy=multi-user.target
SERVICEEOF

echo "==> Service file created at /home/r4wm/hoyj-api.service"

echo "==> Deploying html files..."
sudo cp /home/r4wm/github/hoyj_fs/html/form.html /var/www/html/helpersofyourjoy/mp3/form.html
echo "==> HTML files deployed"
ENDSSH

echo ""
echo "==> Build complete! Binary and service file are ready."
echo ""
echo "Now run these sudo commands on the server:"
echo ""
echo "  ssh $SERVER"
echo "  sudo mv /home/r4wm/hoyj-api.service /etc/systemd/system/"
echo "  sudo systemctl daemon-reload"
echo "  sudo systemctl enable hoyj-api"
echo "  sudo systemctl start hoyj-api"
echo "  sudo systemctl status hoyj-api"
echo ""
echo "Or run this one-liner (will prompt for password):"
echo ""
echo "  ssh -t $SERVER 'sudo mv /home/r4wm/hoyj-api.service /etc/systemd/system/ && sudo systemctl daemon-reload && sudo systemctl enable hoyj-api && sudo systemctl start hoyj-api && sudo systemctl status hoyj-api'"
