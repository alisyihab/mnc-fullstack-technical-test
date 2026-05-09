-- Example data for MNC Fullstack Technical Test - Phase 2
-- Run database migrations first before inserting this data:
--   go run tahap-2/cmd/main.go   (migrations run automatically on startup)
--   OR manually via golang-migrate CLI (see README.md)
--
-- Credentials for testing:
--   User 1: phone=0811255501, PIN=123456
--   User 2: phone=0811255502, PIN=123456

INSERT INTO users (id, first_name, last_name, phone_number, address, pin, balance, created_at, updated_at)
VALUES
    ('bc1c823e-b0fb-4b20-88c0-dff25e283252', 'Guntur', 'Saputro', '0811255501', 'Jl. Kebon Sirih No. 1',
     '$2a$10$MsIOHdljg6PjwH7g1EzfNO7HzvPVqq35u.xauGo7z8/kwkVmu6oyO', 370000,
     '2021-04-01 22:21:20+07', '2021-04-01 22:21:20+07'),
    ('b7342e8e-e8e7-4a5d-873e-b1b1bfcdeddb', 'Tom', 'Araya', '0811255502', 'Jl. Diponegoro No. 215',
     '$2a$10$MsIOHdljg6PjwH7g1EzfNO7HzvPVqq35u.xauGo7z8/kwkVmu6oyO', 30000,
     '2021-04-01 22:20:00+07', '2021-04-01 22:20:00+07');

INSERT INTO transactions (id, user_id, target_user_id, amount, remarks, balance_before, balance_after, status, transaction_type, category, created_at, updated_at)
VALUES
    ('201ddde1-f797-484b-b1a0-07d1190e790a', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', NULL,
     500000, '', 0, 500000, 'SUCCESS', 'CREDIT', 'TOPUP',
     '2021-04-01 22:21:21+07', '2021-04-01 22:21:21+07'),
    ('13bcb11c-111e-4a65-9afd-90a86a01cd21', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', NULL,
     100000, 'Pulsa Telkomsel 100k', 500000, 400000, 'SUCCESS', 'DEBIT', 'PAYMENT',
     '2021-04-01 22:22:00+07', '2021-04-01 22:22:00+07'),
    ('a7d39cf6-44b6-41fc-b3e9-7b16df5321c5', 'bc1c823e-b0fb-4b20-88c0-dff25e283252', 'b7342e8e-e8e7-4a5d-873e-b1b1bfcdeddb',
     30000, 'Hadiah Ultah', 400000, 370000, 'SUCCESS', 'DEBIT', 'TRANSFER',
     '2021-04-01 22:23:20+07', '2021-04-01 22:23:20+07');
