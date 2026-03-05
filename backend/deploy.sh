#!/bin/bash
set -euo pipefail

SERVER="dvolkov@77.73.131.137"
REMOTE_DIR="/opt/teslapay"

echo "Syncing backend to $SERVER:$REMOTE_DIR..."
rsync -avz --exclude='.git' --exclude='*.DS_Store' --exclude='vendor/' \
  /Users/dmitry/dream-team/projects/neobank-teslapay/backend/ \
  $SERVER:$REMOTE_DIR/

echo "Deploying on remote..."
ssh $SERVER "cd $REMOTE_DIR && docker compose pull --ignore-pull-failures && docker compose up -d --build"

echo ""
echo "Deployed! Endpoints:"
echo "  Gateway:  http://77.73.131.137:8080/api/v1"
echo "  Health:   http://77.73.131.137:8080/health"
