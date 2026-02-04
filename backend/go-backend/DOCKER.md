# Docker Setup Guide

This guide explains how to run the Resume Backend using Docker Compose.

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+

## Quick Start

1. **Copy environment file** (if not already present):
   ```bash
   cp .env.example .env
   ```

2. **Update environment variables** in `.env` file:
   - Set `DB_PASSWORD` to a secure password
   - Set `JWT_SECRET` to a strong random secret
   - Configure OAuth providers if needed
   - Configure Cloudinary if using cloud storage

3. **Build and start services**:
   ```bash
   docker-compose up -d
   ```

4. **Check service status**:
   ```bash
   docker-compose ps
   ```

5. **View logs**:
   ```bash
   # All services
   docker-compose logs -f
   
   # Specific service
   docker-compose logs -f backend
   docker-compose logs -f postgres
   docker-compose logs -f migrate
   ```

## Services

### PostgreSQL Database
- **Container**: `resume-postgres`
- **Port**: `5432` (default)
- **Data**: Persisted in `postgres_data` volume
- **Health Check**: Automatically checks database readiness

### Migration Service
- **Container**: `resume-migrate`
- **Purpose**: Runs database migrations before backend starts
- **Runs**: Once on startup, then exits

### Backend API
- **Container**: `resume-backend`
- **Port**: `8080` (default)
- **Uploads**: Persisted in `uploads_data` volume
- **Health Check**: Available at `http://localhost:8080/health`

## Environment Variables

Key environment variables (see `.env.example` for full list):

### Database
- `DB_HOST`: Database host (default: `postgres` in Docker)
- `DB_PORT`: Database port (default: `5432`)
- `DB_USER`: Database user (default: `postgres`)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: `resume_db`)

### Server
- `PORT`: Server port (default: `8080`)
- `UPLOAD_DIR`: Upload directory (default: `/app/uploads` in Docker)

### Authentication
- `JWT_SECRET`: Secret key for JWT tokens (change in production!)
- `JWT_EXPIRY_HOURS`: Access token expiry (default: `24`)
- `REFRESH_EXPIRY_DAYS`: Refresh token expiry (default: `30`)

### OAuth Providers
Configure OAuth providers by setting:
- `GOOGLE_ENABLED=true`
- `GOOGLE_CLIENT_ID=your-client-id`
- `GOOGLE_CLIENT_SECRET=your-client-secret`

Same pattern for `GITHUB_*` and `MICROSOFT_*` variables.

## Common Commands

### Start services
```bash
docker-compose up -d
```

### Stop services
```bash
docker-compose down
```

### Stop and remove volumes (⚠️ deletes data)
```bash
docker-compose down -v
```

### Rebuild services
```bash
docker-compose build
docker-compose up -d
```

### Run migrations manually
```bash
docker-compose run --rm migrate up
```

### Check migration status
```bash
docker-compose run --rm migrate -command=status
```

### Access database
```bash
docker-compose exec postgres psql -U postgres -d resume_db
```

### Access backend container shell
```bash
docker-compose exec backend sh
```

### View backend logs
```bash
docker-compose logs -f backend
```

## Development

For local development with hot reload, you can use `docker-compose.override.yml`:

1. Copy the example:
   ```bash
   cp docker-compose.override.yml.example docker-compose.override.yml
   ```

2. Modify as needed (e.g., mount source code for live editing)

## Troubleshooting

### Database connection issues
- Ensure PostgreSQL is healthy: `docker-compose ps postgres`
- Check database logs: `docker-compose logs postgres`
- Verify environment variables: `docker-compose exec backend env | grep DB_`

### Migration failures
- Check migration logs: `docker-compose logs migrate`
- Run migrations manually: `docker-compose run --rm migrate up`
- Check migration status: `docker-compose run --rm migrate -command=status`

### Backend not starting
- Check logs: `docker-compose logs backend`
- Verify health endpoint: `curl http://localhost:8080/health`
- Check if port is in use: `netstat -an | grep 8080` (Linux/Mac) or `netstat -an | findstr 8080` (Windows)

### OCR/PDF processing issues
- Ensure Tesseract is installed in container: `docker-compose exec backend tesseract --version`
- Ensure Poppler is installed: `docker-compose exec backend pdftoppm -v`
- Check upload directory permissions: `docker-compose exec backend ls -la /app/uploads`

## Production Considerations

1. **Change default passwords**: Update `DB_PASSWORD` and `JWT_SECRET`
2. **Use secrets management**: Consider Docker secrets or external secret managers
3. **Enable SSL**: Set `DB_SSLMODE=require` for database connections
4. **Resource limits**: Add resource constraints in `docker-compose.yml`
5. **Backup strategy**: Implement regular backups for `postgres_data` volume
6. **Monitoring**: Add monitoring and logging solutions
7. **Reverse proxy**: Use nginx/traefik for SSL termination and routing

## Volumes

- `postgres_data`: PostgreSQL database files
- `uploads_data`: Uploaded files (if using local storage)

To backup volumes:
```bash
docker run --rm -v resume-backend_postgres_data:/data -v $(pwd):/backup alpine tar czf /backup/postgres_backup.tar.gz /data
```
