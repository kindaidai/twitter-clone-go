#!/bin/bash
# -xで実行コマンドを表示、-eでエラーがあったときに終了
set -xe
set -o pipefail

export MYSQL_PORT=${MYSQL_PORT:-3306}
export MYSQL_USER=${MYSQL_USER:-user}
export MYSQL_DBNAME=${MYSQL_DBNAME:-sample}
export MYSQL_PWD=${MYSQL_PASS:-password}
export LANG="C.UTF-8"

# connection test
mysql -h localhost -P $MYSQL_PORT -u $MYSQL_USER $MYSQL_DBNAME
