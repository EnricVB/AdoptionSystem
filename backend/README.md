# AdoptionSystem - Backend

Proyecto Backend del sistema `AdoptionSystem` desarrollado con Go.
Este sistema contiene los recursos necesarios para realizar peticiones `APIRest` para tener una conexión segura a `BBDD`, así como integración de sistemas de Authentication con `Google OAuth`, `encriptación` de contraseñas, y sistema de `migraciones`.
Además, soporta la carga de imagenes de las entradas de la BBDD para un almacenamiento más optimo con unas caracteristicas similares en todas las imagenes.

## Setup
Para preparar el Backend para su uso se han de modificar el fichero .env, el cual, actualmente contiene datos de prueba.
Se utiliza la librería `Godotenv` para cargar como variables de entorno los datos de `.env`.

`GoDotEnv` neceista que el `.env` se localice en la misma carpeta que el `main.go`, tenlo en cuenta si vas a modificar las carpetas.

Tras tener las .env configuradas, deberás lanzar tu MySQL, con los datos del `.env`, y finalmente solo te falta realizar pruebas.

Por último, ejecuta el proyecto con `go run main.go` desde la carpeta `backend/`.

⚠️ `Si utilizas el docker de pruebas, recuerda ejecutar 'docker start mysql-container' para iniciar la BBDD` 


## API REST

El sistema APIRest acepta peticiones únicamente desde el puerto `:4200` o desde el propio puerto del backend `:8080`, con la arquitectura en capas `Route -> Handler -> Service -> DAO`, donde:

- Route:
Define el endpoint y enlaza la ruta HTTP con un handler concreto mediante la librería `Echo` y realiza unas validaciones iniciales.

- Handler:
Recibe los datos de la petición, valida datos de entrada y llama al `servicio` correspondiente para realizar todas las operaciones internas.
No contiene lógica de negocio, solo valida que los datos sean correctos, no causen excepciones, etc.

- Service:
Contiene la lógica de negocio, llamando al `DAO` o realizando las operaciones necesarias de persistencia y peticiones de otros datos a otros DAO.

- DAO:
Se encarga de realizar todas las peticiones internas a BBDD, de forma que se encapsulan las peticiones con las validaciones necesarias.

### Testing
Este sistema de API Rest soporta testing manual y automatico mediante `rest_client.http`, el cual permite, mediante parametros, realizar desde el propio IDE `VSCode + REST Client`, peticiones postman robotizadas hacia el sistema de APIRest.  


## Mail

Sistema de Mailing, utilizando la libreria Go-Mail, permite enviar mail automatizados y manuales entre varios clientes y el administrador.

De forma que se abstrae la lógica de mailing a llamadas con la API Rest con unos parametros como `Cliente, Motivo`.
Dependiendo el endpoint utilizado, se utilizará un template en concreto, pudiendo así automatizar los envios sin sobrecargar el código de HTML.

## Marshal - Error

Existe un sistema propio para el manejo de respuestas HTTP mediante las clases de utilidad `error` y `marshal`, así pudiendo formatear las respuestas tanto de error como las correctas con un mismo formato.

### Error

Utilizando esta utilidad, solo se ha de enviar por parametro el código de error y el mensaje, formateando la respuesta con:

{
    code: YYY,
    message: "XXXX"
}


