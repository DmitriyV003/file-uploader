CREATE TABLE uploaded_files
(
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    ext            VARCHAR(50)  NOT NULL,
    original_name  VARCHAR(255) NOT NULL,
    generated_name VARCHAR(255) NOT NULL,
    created_at     DATETIME     NOT NULL,
    updated_at     DATETIME     NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE servers
(
    id   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    url VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE uploaded_file_chunks
(
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    uploaded_file_id BIGINT UNSIGNED NOT NULL,
    chunk_number     INT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    size             BIGINT UNSIGNED NOT NULL,
    server_id        BIGINT UNSIGNED NOT NULL,
    hash             VARCHAR(255),
    created_at       DATETIME     NOT NULL,
    updated_at       DATETIME     NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_uploaded_file
        FOREIGN KEY (uploaded_file_id) REFERENCES uploaded_files (id),
    CONSTRAINT fk_server
        FOREIGN KEY (server_id) REFERENCES servers (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `servers` (url) VALUES ('store-server-base:50052');
INSERT INTO `servers` (url) VALUES ('store-server-1:50053');
INSERT INTO `servers` (url) VALUES ('store-server-2:50054');
INSERT INTO `servers` (url) VALUES ('store-server-3:50055');
INSERT INTO `servers` (url) VALUES ('store-server-4:50056');
INSERT INTO `servers` (url) VALUES ('store-server-5:50057');
