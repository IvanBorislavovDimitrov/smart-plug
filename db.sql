CREATE TABLE plugs(
	id VARCHAR(255) PRIMARY KEY,
	name VARCHAR(255) UNIQUE NOT NULL,
	ip_address VARCHAR(255) UNIQUE NOT NULL,
	power_to_turn_off DECIMAL NOT NULL,
	created_at TIMESTAMP NOT NULL
)