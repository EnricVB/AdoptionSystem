ALTER TABLE Pets
  ADD COLUMN Genre ENUM('Male', 'Female'),
  ADD COLUMN Weight DECIMAL(5,2) NOT NULL DEFAULT 0.00;

CREATE TABLE VaccinationHistory (
  ID INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  Pet_ID INT UNSIGNED NOT NULL,
  Vaccination_Date DATE NOT NULL,
  Vaccine_Name VARCHAR(100) NOT NULL,
  FOREIGN KEY (Pet_ID) REFERENCES Pets(ID)
);

-- Para ejecutar este archivo SQL con Go y una herramienta de migraciones como golang-migrate, usa el comando:
-- migrate -path /root/adoption-system/backend/cmd/migrations -database "tu_cadena_de_conexion" up