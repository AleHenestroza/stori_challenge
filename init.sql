-- Create database
CREATE DATABASE stori;
-- Create role
CREATE ROLE stori WITH LOGIN PASSWORD 'stori' CREATEDB;
-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE stori TO stori;
-- Connect to database
\c stori
-- Create tables
CREATE TABLE IF NOT EXISTS transactions (
    Id SERIAL PRIMARY KEY,
    Date DATE NOT NULL,
    Amount DECIMAL NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    Email VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS user_transactions (
    Id_User INT REFERENCES users(Id),
    Id_transaction INT REFERENCES transactions(Id),
    PRIMARY KEY (Id_User, Id_transaction)
);
GRANT SELECT ON transactions TO PUBLIC;
INSERT INTO transactions (Id, Date, Amount)
SELECT
    data.id,
    TO_TIMESTAMP(data.date, 'YYYY/MM/DD')::DATE AS date,
    data.amount::DECIMAL
FROM (
    VALUES
        (0,'2022/7/15',60.5),
        (1,'2022/7/18',-10.3),
        (2,'2022/8/2',-20.46),
        (3,'2022/8/13',10),
        (4,'2022/8/20',-5.75),
        (5,'2022/9/5',25.8),
        (6,'2022/9/12',-12.2),
        (7,'2022/9/27',30.1),
        (8,'2022/10/5',-15.75),
        (9,'2022/10/18',18.9),
        (10,'2022/10/25',-7.6),
        (11,'2022/11/7',12.3),
        (12,'2022/11/20',-9.8),
        (13,'2022/12/3',14.6),
        (14,'2022/12/16',-18.2),
        (15,'2023/1/2',22.4),
        (16,'2023/1/15',-5.9),
        (17,'2023/1/28',8.7),
        (18,'2023/2/10',-16.3),
        (19,'2023/2/23',7.2),
        (20,'2023/3/8',-11.1),
        (21,'2023/3/21',20.5),
        (22,'2023/4/3',-6.4),
        (23,'2023/4/16',16.9),
        (24,'2023/4/29',-13.7),
        (25,'2023/5/12',9.8),
        (26,'2023/5/25',-10.2),
        (27,'2023/6/7',15.3),
        (28,'2023/6/20',-8.1),
        (29,'2023/6/30',12.7)
) AS data(id, date, amount)
WHERE NOT EXISTS (
    SELECT 1
    FROM transactions
    LIMIT 1
);