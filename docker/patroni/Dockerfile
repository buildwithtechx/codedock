FROM postgres:16-bookworm

RUN apt-get update && apt-get install -y --no-install-recommends python3 python3-pip && rm -rf /var/lib/apt/lists/* \
    && pip3 install --no-cache-dir --break-system-packages "patroni[etcd3]==4.0.4" "psycopg2-binary==2.9.10"

COPY patroni.yml /etc/patroni.yml
COPY post_bootstrap.py /post_bootstrap.py
COPY entrypoint.sh /entrypoint.sh
RUN chmod 644 /etc/patroni.yml && chmod 755 /post_bootstrap.py /entrypoint.sh

EXPOSE 5432 8008

ENTRYPOINT ["/entrypoint.sh"]
