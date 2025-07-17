user=secret_user
password=secret_password
port=5672
ui_port=15672
volume_path="$PWD/rabbitmq_data"

if [ ! -d "$volume_path" ]; then
    mkdir "$volume_path"
    echo "MQ Folder created"
fi

docker run -d \
  -p $port:5672 \
  -p $ui_port:15672 \
  -e RABBITMQ_DEFAULT_USER=$user \
  -e RABBITMQ_DEFAULT_PASS=$password \
  -v $volume_path:/var/lib/rabbitmq \
  rabbitmq:4-management-alpine