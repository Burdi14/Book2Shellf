# Book2Shellf

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

## Admin Panel

Access the admin panel at: `/book2shadmin`

**Default credentials:**
- Username: `admin`
- Password: `changeme`

## API

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
eMIT
