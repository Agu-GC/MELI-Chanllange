CREATE TABLE databases (
    database_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    connection_details JSON,
    analyzed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY (name, server_name)
);

CREATE TABLE tables (
    table_id INT PRIMARY KEY AUTO_INCREMENT,
    database_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    row_count INT,
    description TEXT,
    FOREIGN KEY (database_id) REFERENCES databases(database_id),
    UNIQUE KEY (database_id, name)
);

CREATE TABLE columns (
    column_id INT PRIMARY KEY AUTO_INCREMENT,
    table_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    data_type VARCHAR(100) NOT NULL,
    is_nullable BOOLEAN,
    confidence DECIMAL(5,4) CHECK (confidence BETWEEN 0 AND 1),
    classified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (table_id) REFERENCES tables(table_id),
    UNIQUE KEY (table_id, name)
);

CREATE TABLE classification_rules (
    rule_id INT PRIMARY KEY AUTO_INCREMENT,
    pattern VARCHAR(255),
    data_type VARCHAR(100)
    last_modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
);