ALTER TABLE Pets
  ADD COLUMN Genre ENUM('Male', 'Female'),
  ADD COLUMN Weight DECIMAL(5,2) NOT NULL DEFAULT 0.00;

CREATE TABLE VaccinationHistory (
  Pet_ID INT UNSIGNED NOT NULL,
  Vaccination_Date DATE NOT NULL,
  Vaccine_Name VARCHAR(100) NOT NULL,
  PRIMARY KEY (Pet_ID, Vaccination_Date, Vaccine_Name),
  FOREIGN KEY (Pet_ID) REFERENCES Pets(ID)
);

-- Para ejecutar este archivo SQL con Go y una herramienta de migraciones como golang-migrate, usa el comando:
-- migrate -path /root/adoption-system/backend/cmd/migrations -database "tu_cadena_de_conexion" up