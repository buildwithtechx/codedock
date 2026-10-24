#!/usr/bin/env python3
import os

import psycopg2
from psycopg2 import sql

password = os.environ["CODEDOCK_PG_PASSWORD"]
conn = psycopg2.connect(host="/tmp", user="postgres", dbname="postgres")
conn.autocommit = True
cur = conn.cursor()
cur.execute(
    sql.SQL("CREATE ROLE codedock WITH LOGIN PASSWORD {} CREATEDB").format(
        sql.Literal(password)
    )
)
cur.execute(sql.SQL("CREATE DATABASE codedock OWNER codedock"))
cur.close()
conn.close()
