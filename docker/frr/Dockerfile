FROM quay.io/frrouting/frr:9.1.3

COPY docker-start.sh /usr/local/bin/frr-start.sh
RUN chmod +x /usr/local/bin/frr-start.sh

ENTRYPOINT ["/usr/local/bin/frr-start.sh"]
