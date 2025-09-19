CREATE TABLE rule_attributes (
    id SERIAL PRIMARY KEY,
    name varchar(50) NOT NULL UNIQUE,
    description varchar(300)
);

INSERT INTO rule_attributes (name, description) VALUES
('user.id', 'Unique user identifier'),
('user.email', 'User email address'),
('user.anonymous', 'Is user anonymous'),
('country_code', 'ISO country code'),
('region', 'Region or state'),
('city', 'City name'),
('manufacturer', 'Device manufacturer'),
('device_type', 'Device type (mobile, desktop, tablet, etc.)'),
('os', 'Operating system'),
('os_version', 'Operating system version'),
('browser', 'Browser name'),
('browser_version', 'Browser version'),
('language', 'User language/locale'),
('connection_type', 'Network connection type'),
('age', 'User age'),
('gender', 'User gender'),
('ip', 'IP address'),
('app_version', 'Application version'),
('platform', 'Platform (web, android, ios)');
