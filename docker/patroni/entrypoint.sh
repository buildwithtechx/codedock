#!/bin/sh
set -e
mkdir -p /var/lib/postgresql/data
chown -R postgres:postgres /var/lib/postgresql/data
exec gosu postgres patroni /etc/patroni.yml
