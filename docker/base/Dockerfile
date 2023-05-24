FROM debian:bullseye-slim

COPY --from=cribl/scope-ebpf-builder:local /bin/scope-ebpf /usr/local/bin/scope-ebpf

CMD ["/usr/local/bin/scope-ebpf"]
