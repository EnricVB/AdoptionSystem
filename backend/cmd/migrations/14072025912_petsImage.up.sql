ALTER TABLE ADOPTION_SYS.Pets
  CHANGE COLUMN Is_Adopted Status ENUM('Available', 'FosterHome', 'Adopted') NOT NULL DEFAULT 'Available',
  ADD COLUMN Vaccinated BOOLEAN NOT NULL DEFAULT FALSE;

-- Para ejecutar este archivo SQL con Go y una herramienta de migraciones como golang-migrate, usa el comando:
-- migrate -path /root/adoption-system/backend/cmd/migrations -database "tu_cadena_de_conexion" up