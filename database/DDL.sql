CREATE DATABASE IF NOT EXISTS xyz_multifinance;
USE xyz_multifinance;

CREATE TABLE customers (
    nik VARCHAR(16) PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    birth_place VARCHAR(100) NOT NULL,
    birth_date DATE NOT NULL,
    salary BIGINT NOT NULL,
    photo_ktp VARCHAR(255) NOT NULL,
    photo_selfie VARCHAR(255) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE customer_limits (
    customer_nik VARCHAR(16),
    tenor INT NOT NULL,
    limit_amount BIGINT NOT NULL,
    PRIMARY KEY (customer_nik, tenor),
    FOREIGN KEY (customer_nik) REFERENCES customers(nik)
) ENGINE=InnoDB;

CREATE TABLE transactions (
    contract_number VARCHAR(50) PRIMARY KEY,
    customer_nik VARCHAR(16),
    otr BIGINT NOT NULL,
    admin_fee BIGINT NOT NULL,
    installment BIGINT NOT NULL,
    interest BIGINT NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL,
    FOREIGN KEY (customer_nik) REFERENCES customers(nik)
) ENGINE=InnoDB;

