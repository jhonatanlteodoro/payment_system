password=secret_password
port=6379
volume_path="$PWD/redis_data"

if [ ! -d "$volume_path" ]; then
    mkdir "$volume_path"
    echo "Redis Folder created"
fi

docker run -d \
  -p $port:6379 \
  -v $volume_path:/data \
  redis:8-alpine \
  redis-server --requirepass $password --appendonly yes