ALTER TABLE Users
  ADD COLUMN Provider VARCHAR(50) NOT NULL DEFAULT 'local',
  ADD COLUMN Provider_ID VARCHAR(100) UNIQUE;

-- Para ejecutar este archivo SQL con Go y una herramienta de migraciones como golang-migrate, usa el comando:
-- migrate -path /root/adoption-system/backend/cmd/migrations -database "tu_cadena_de_conexion" up