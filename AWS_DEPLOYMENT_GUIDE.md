# TeamSync AWS Deployment & Debugging Guide

This document summarizes the steps taken to deploy TeamSync on an AWS EC2 `t3.micro` instance running Ubuntu 24.04.

## 1. Infrastructure Setup (AWS Console)
- **AMI:** Ubuntu 24.04 LTS (64-bit x86)
- **Instance Type:** `t3.micro` (1GB RAM)
- **Storage:** 20GB EBS (gp3)
- **Security Group Rules:**
    - SSH (22): For remote access.
    - HTTP (80): For the web frontend.
    - (Optional) Port 8080: If direct backend access is needed.

## 2. Server Initialization (Run on EC2)
The following script handles Docker installation and memory optimization (Swap).

```bash
# Update and install Docker
sudo apt-get update -y
sudo apt-get install -y docker.io docker-compose
sudo systemctl enable --now docker
sudo usermod -aG docker ubuntu

# Add 2GB Swap Space (CRITICAL for t3.micro to prevent build crashes)
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## 3. Environment Configuration
Created at `~/teamsync/.env` on the server:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=<random_generated>
DB_NAME=teamsync
REDIS_ADDR=redis:6379
JWT_SECRET=<random_generated>
FROM_MAIL=<email>
PASS_MAIL=<password>
PORT=8080
```

## 4. Deployment Commands

### Local to Remote Transfer (Run on Mac)
```bash
chmod 400 ~/Downloads/user.pem
rsync -avz -e "ssh -i ~/Downloads/user.pem" \
  --exclude 'node_modules' \
  --exclude '.git' \
  --exclude 'frontend/node_modules' \
  ./ ubuntu@13.63.66.118:~/teamsync/
```

### Domain Configuration
- **Domain:** teamsynch.tech
- **Elastic IP:** 13.63.66.118
- **A Record:** Point `@` to `13.63.66.118`
- **CNAME Record:** Point `www` to `teamsynch.tech`

### Starting the App (Run on EC2)
```bash
cd ~/teamsync
sudo docker-compose up -d --build
```

## 5. Debugging & Maintenance

| Command | Purpose |
|---------|---------|
| `sudo docker ps` | Check if all 4 containers are running |
| `sudo docker-compose logs -f` | Stream logs from all services |
| `sudo docker logs teamsync-backend-1` | Debug backend/email issues specifically |
| `sudo docker-compose down` | Stop and remove all containers |
| `sudo docker-compose exec db psql -U user teamsync` | Enter the Postgres database CLI |
| `free -m` | Check memory and swap usage |

## 6. Common Issues
- **Connection Refused:** Ensure the Security Group allows Port 80.
- **Build Crashes:** Usually due to memory. Check `free -m` to ensure Swap is active.
- **Emails not sending:** Check the `PASS_MAIL` (App Password) and ensure the backend logs don't show "Authentication Failed".
- **Database not ready:** The backend has a `depends_on` healthcheck; if it won't start, check `docker logs teamsync-db-1`.
