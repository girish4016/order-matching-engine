# order-matching-engine
mimics stock exchange order matching system


# 1. Start TiDB (in a separate terminal or use tmux/screen)
tiup playground

# --- Open a new terminal for the following commands ---

# 2. Connect to TiDB using MySQL CLI
mysql -h 127.0.0.1 -P 4000 -u root

# 3. Inside the MySQL CLI, paste the following SQL commands:

CREATE DATABASE IF NOT EXISTS order_matching_engine;
USE order_matching_engine;

CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(36) PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    side ENUM('buy', 'sell') NOT NULL,
    type ENUM('limit', 'market') NOT NULL,
    price DECIMAL(18,8) NOT NULL,
    quantity DECIMAL(18,8) NOT NULL,
    remaining_quantity DECIMAL(18,8) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE trades (
    id VARCHAR(36) PRIMARY KEY,
    buy_order_id VARCHAR(36) NOT NULL,
    sell_order_id VARCHAR(36) NOT NULL,
    price DECIMAL(20,8) NOT NULL,
    quantity DECIMAL(20,8) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    symbol VARCHAR(20) NOT NULL
);

# 4. Exit MySQL CLI
exit;

# 5. Run your Go server (adjust path if needed)
go run cmd/server/main.go
