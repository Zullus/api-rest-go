CREATE DATABASE IF NOT EXISTS apirestgo DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE apirestgo;

CREATE TABLE currency_exchange_log (
  UUID varchar(38) NOT NULL,
  amount decimal(10,2) NOT NULL,
  converted_amount decimal(10,2) NOT NULL,
  from_currency varchar(5) NOT NULL,
  to_currency varchar(5) NOT NULL,
  conversion_rate decimal(10,2) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE currency_exchange_log
  ADD PRIMARY KEY (UUID);
COMMIT;
