CREATE TABLE schema_company (
  ID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  Company_Name VARCHAR(255) NOT NULL,
  Contact_Email VARCHAR(255) NOT NULL,
  Phone VARCHAR(20),
  Address VARCHAR(500),
  City VARCHAR(100),
  State VARCHAR(100),
  Postal_Code VARCHAR(20),
  Country VARCHAR(100),
  Created_At TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Updated_At TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Para ejecutar este archivo SQL con Go y una herramienta de migraciones como golang-migrate, usa el comando:
-- migrate -path /root/adoption-system/backend/cmd/migrations -database "tu_cadena_de_conexion" up