CREATE TABLE IF NOT EXISTS classifications (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    pattern TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    sensitivity_level INT NOT NULL,
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO classifications (name, description, pattern, category, sensitivity_level) VALUES
('EMAIL', 'Emails Columns',
 '(?i)^(?:user_|client_|person_)?(?:e[-]?mail|mail_?addr(?:ess)?|correo_?electronico?|email_?(?:id|name))$', 'PII', 3),

('HEALTH_DATA', 'Health data columns',
 '(?i)^(?:medical_|health_)(?:record|history|condition|bmi|blood_?type|diagnosis|treatment|prescription)\b', 'Healthcare', 4),

('PASSWORD', 'Password columns',
 '(?i)^(?:user_|account_)?(?:password|pwd|pass_?(?:hash|phrase)|secret_?(?:key|token)|auth_?token|security_?key)\b', 'Security', 5),

('CREDIT_CARD', 'Credit cards columns',
 '(?i)^(?:(?:cc_|credit_?card_?)(?:number|num|no|hash|token|info|data|holder_?name|expiration_?(?:month|year)|cvv)|(?:card_holder_?name|card_?number|expiration_?(?:month|year)|cvv))$', 'Financial', 5),

('PHONE_NUMBER', 'Phone number columns',
 '(?i)^(?:phone|mobile|tel|cell|contact_?no)(?:_*(?:number|num|code|country|detail))?$', 'Contact Data', 3),

('IDENTIFICATION', 'ID columns',
 '(?i)^(?:national_?|gov_?|state_?|country_?)?(?:id|identification|ssn|tax_?id|nin|curp|dl|driver_?license|passport)$', 'Government ID', 4),

('ADDRESS', 'Address columns',
 '(?i)^(?:postal_|physical_|shipping_|billing_)?(?:address|direccion|adresse|addr_?line|street|residence|city|zip_?code)\b', 'PII', 2),

('FINANCIAL', 'Financial information columns',
 '(?i)^(?:annual_|monthly_|daily_)?(?:salary|income|revenue|balance|tax|credit_?score|expense|payment|invoice)\b', 'Financial', 4),

('GEOLOCATION', 'GPS data columns',
 '(?i)^(?:gps_|geo_|location_?)(?:coordinates|lat_?long|position_?data|latitude|longitude|altitude)\b', 'Location Data', 3);