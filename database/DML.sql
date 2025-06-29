INSERT INTO customers (nik, full_name, legal_name, birth_place, birth_date, salary, photo_ktp, photo_selfie)
VALUES 
    ('1234567890123456', 'Budi', 'Budi Santoso', 'Jakarta', '1990-01-01', 10000000, 'ktp_budi.jpg', 'selfie_budi.jpg'),
    ('9876543210987654', 'Annisa', 'Annisa Rahma', 'Bandung', '1995-05-05', 15000000, 'ktp_annisa.jpg', 'selfie_annisa.jpg');

INSERT INTO customer_limits (customer_nik, tenor, limit_amount)
VALUES 
    ('1234567890123456', 1, 100000),
    ('1234567890123456', 2, 200000),
    ('1234567890123456', 3, 500000),
    ('1234567890123456', 4, 700000),
    ('9876543210987654', 1, 1000000),
    ('9876543210987654', 2, 1200000),
    ('9876543210987654', 3, 1500000),
    ('9876543210987654', 4, 2000000);