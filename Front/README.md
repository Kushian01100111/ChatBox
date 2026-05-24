# ChatBox Vue Client

Cliente simple en Vue/Vite para interactuar con el backend Go del archivo `Back.rar`.

La API del backend no expone endpoints REST para mensajes. La comunicación principal es por WebSocket:

```txt
ws://localhost:8080/ws/:roomId
```

## Requisitos

- Node.js instalado
- Backend Go corriendo en `localhost:8080`

## Correr el backend

Desde la carpeta del backend:

```bash
go run ./cmd/api
```

Si tu Go local no soporta la versión indicada en `go.mod`, podés cambiar temporalmente:

```txt
go 1.25.0
```

por una versión instalada en tu máquina, por ejemplo:

```txt
go 1.23.0
```

## Correr el frontend

Desde esta carpeta:

```bash
npm install
npm run dev
```

Abrí la URL que te muestre Vite, normalmente:

```txt
http://localhost:5173
```

## Probar chat de sala

1. Abrí dos pestañas del frontend.
2. En ambas conectate a:

```txt
room-1
```

3. En una pestaña enviá un mensaje de sala.
4. La otra pestaña debería recibirlo.

El payload que manda la app es equivalente a:

```json
{
  "type": "message",
  "id": "room-1",
  "sender": "Pedro",
  "content": "Hola desde Vue"
}
```

## Probar notificaciones

1. Abrí una pestaña conectada como:

```txt
user-123
```

2. Abrí otra pestaña conectada como:

```txt
room-1
```

3. Desde la pestaña `room-1`, enviá una notificación a:

```txt
user-123
```

El payload es equivalente a:

```json
{
  "type": "notification",
  "sender": "Pedro",
  "recipient": "user-123",
  "content": "Tenés una nueva notificación"
}
```

## Nota de diseño

El backend actual usa el parámetro `/ws/:roomId` tanto para salas como para usuarios. Para una versión más seria convendría separar `roomId` y `userId`, pero para probar la API actual esta app funciona bien.
