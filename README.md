## Tech Stack

- **Backend**: Go (Gin framework)
- **Frontend**: React
- **Database**: SQLite
- **Styling**: CSS 

## Getting Started

```bash
docker-compose up -d

open http://localhost:8000
```

### Manual Setup

**Backend:**
```bash
cd app
go mod tidy
go run ./main
```

**Frontend:**
```bash
cd app/frontend
npm install
npm start
```

## Admin Access

Access the admin panel at: `/book2shadmin`

**Default credentials:**
- Username: `admin`
- Password: ``

> ⚠️ **Change these in production** by setting environment variables:
> - `ADMIN_USERNAME`
> - `ADMIN_PASSWORD`

## Project Structure

```
├── app/
│   ├── handlers/         # Go handlers and database
│   │   ├── auth.go       # Authentication
│   │   ├── database.go   # SQLite operations
│   │   ├── handlers.go   # API endpoints
│   │   └── models.go     # Data models
│   ├── main/
│   │   └── main.go       # Entry point
│   └── frontend/         # React application
│       ├── src/
│       │   ├── components/
│       │   ├── pages/
│       │   ├── api.js
│       │   └── index.css  # Hacker theme styles
│       └── public/
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Public
- `GET /api/books`
- `GET /api/books/:id`
- `GET /api/books/:id/download`
- `GET /api/sections`
- `GET /api/sections/:id/books`

### Admin (Protected)
- `POST /api/login`
- `POST /api/admin/books`
- `PUT /api/admin/books/:id`
- `DELETE /api/admin/books/:id`
- `POST /api/admin/sections`
- `PUT /api/admin/sections/:id`
- `DELETE /api/admin/sections/:id`
- `POST /api/admin/upload/book`
- `POST /api/admin/upload/cover`

## License
MIT
