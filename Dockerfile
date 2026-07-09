# Stage 1: Pull from official RustFS (or use a specific tag)
FROM rustfs/rustfs:latest AS rustfs-source
FROM rustfs/rc:latest AS rc-source
FROM minio/minio:latest AS minio-source

# Stage 2: Hummingbird base (choose the appropriate minimal image)
FROM registry.access.redhat.com/hi/core-runtime:latest-builder


# Copy the main binary and any required files
COPY --from=rustfs-source /usr/bin/rustfs /usr/bin/rustfs
COPY --from=rustfs-source /entrypoint.sh /entrypoint.sh
COPY --from=rc-source /usr/bin/rc /usr/bin/rc
COPY --from=minio-source /usr/bin/mc /usr/bin/mc

USER root
# Set permissions, user (Hummingbird often prefers non-root)
RUN mkdir -p /opt/rustfs/events && \
    chmod +x /usr/bin/rustfs /usr/bin/rc /entrypoint.sh && \
    useradd -r -u 1001 -m rustfs && \
    mkdir -p /data /logs && \
    chown -R rustfs:rustfs /data /logs /opt/rustfs

USER 1001

ENV RUSTFS_CONSOLE_CORS_ALLOWED_ORIGINS="*" \
    RUSTFS_VOLUMES="/data" \
    RUSTFS_OBS_LOGGER_LEVEL=warn \
    RUSTFS_OBS_LOG_DIRECTORY= \
    RUSTFS_OBS_ENVIRONMENT=production

ENV RUSTFS_AUDIT_ENABLE=true \
    RUSTFS_AUDIT_WEBHOOK_ENABLE_PRIMARY="on" \
    RUSTFS_AUDIT_WEBHOOK_ENDPOINT_PRIMARY="https://localhost:9999/"

EXPOSE 9000 9001

VOLUME ["/data"]

ENTRYPOINT ["/entrypoint.sh"]

CMD ["rustfs"]
