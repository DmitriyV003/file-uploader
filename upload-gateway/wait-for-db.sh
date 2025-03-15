#!/bin/sh

set -e

host="$1"
shift
cmd="$@"

until nc -z "$host" 3306; do
  >&2 echo "Ожидание базы данных на $host:3306..."
  sleep 1
done

>&2 echo "База данных доступна, запускаем миграции..."
soda migrate

exec $cmd
