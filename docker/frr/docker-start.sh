#!/bin/sh
set -e

install -d -o frr -g frr -m 0755 /etc/frr
install -d -o frr -g frr -m 0755 /var/run/frr

if [ ! -f /etc/frr/frr.conf ]; then
    echo "frr version 9.1.3" > /etc/frr/frr.conf
    echo "frr defaults traditional" >> /etc/frr/frr.conf
    echo "!" >> /etc/frr/frr.conf
fi

ulimit -n 65536

exec /usr/lib/frr/bgpd \
    --no_kernel \
    -f /etc/frr/frr.conf \
    -i /var/run/frr/bgpd.pid \
    --vty_socket /var/run/frr
