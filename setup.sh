#!/bin/bash
if [ -z "$1" ]; then
    echo "Error: AES_KEY is required"
    exit 1
fi
AES_KEY=$1
if [ ! -d "$HOME/tao-kieu-chu/assets" ]; then
    echo "Warning: assets directory does not exist"
    exit 1
fi

if [ ! -d "$HOME/tao-kieu-chu/fonts" ]; then
    echo "Warning: fonts directory does not exist"
    exit 1
fi

if [ ! -d "$HOME/tao-kieu-chu/views" ]; then
    echo "Warning: views directory does not exist"
    exit 1
fi
podman pull ghcr.io/43vn/tao-kieu-chu:latest
if [ $? -ne 0 ]; then
    echo "Error: Failed to pull image"
    exit 1
fi
if podman ps -a --format '{{.Names}}' | grep -q "tao-kieu-chu"; then
    podman rm -f tao-kieu-chu
fi
podman run -d --name tao-kieu-chu \
    -e PORT=28182 \
    -e AES_KEY="$AES_KEY" \
    -v "$HOME/tao-kieu-chu/assets:/app/assets" \
    -v "$HOME/tao-kieu-chu/fonts:/app/fonts" \
    -v "$HOME/tao-kieu-chu/views:/app/views" \
    --restart=always \
    --network=host \
    ghcr.io/43vn/tao-kieu-chu:latest
echo "Container 'tao-kieu-chu' started successfully."
podman logs tao-kieu-chu
