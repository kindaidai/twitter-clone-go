#!/bin/bash
# -xで実行コマンドを表示、-eでエラーがあったときに終了
set -xe

CURRENT_DIR=$(cd $(dirname $0);pwd)
export MYSQL_HOST=${MYSQL_HOST:-localhost}
export MYSQL_PORT=${MYSQL_PORT:-3306}
export MYSQL_USER=${MYSQL_USER:-sample_user}
export MYSQL_DBNAME=${MYSQL_DBNAME:-sample}
export MYSQL_PWD=${MYSQL_PASSWORD:-password}
export LANG="C.UTF-8"
cd $CURRENT_DIR

printenv

# connection test
mysql --defaults-file=/dev/null -h $MYSQL_HOST -P $MYSQL_PORT -u $MYSQL_USER $MYSQL_DBNAME
