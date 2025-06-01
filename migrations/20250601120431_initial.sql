-- +goose Up
-- +goose StatementBegin
CREATE TABLE `users` (
    `id` bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(3) NULL,
    `updated_at` datetime(3) NULL,
    `deleted_at` datetime(3) NULL,
    `login` varchar(16) NOT NULL,
    `password` varchar(255) NOT NULL,
    `twilio_account_s_id` varchar(34) NOT NULL,
    `twilio_auth_token` varchar(255) NOT NULL,
    `callback_uuid` varchar(36) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `idx_users_deleted_at` (`deleted_at`),
    CONSTRAINT `uni_users_login` UNIQUE (`login`),
    CONSTRAINT `uni_users_callback_uuid` UNIQUE (`callback_uuid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
-- +goose StatementEnd
---
-- +goose Down
-- +goose StatementBegin
DROP TABLE `users`;
-- +goose StatementEnd