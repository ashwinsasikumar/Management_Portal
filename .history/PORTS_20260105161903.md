# Management Portal - Port Configuration

## Standard Ports

This application uses the following standard ports:

| Service   | Host Port | Container Port | Description                          |
|-----------|-----------|----------------|--------------------------------------|
| Frontend  | 8082      | 80             | React application served by Nginx    |
| Backend   | 8080      | 8080           | Go API server                        |
| MySQL     | 3308      | 3306           | MySQL database                       |

## Access URLs

- **Frontend Application**: http://localhost:8082
- **Backend API**: http://localhost:8080
- **MySQL Database**: localhost:3308

## API Configuration

When configuring the frontend to communicate with the backend:

```javascript
// In React app
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
```

## Database Connection

When connecting to MySQL from outside Docker:

```bash
mysql -h 127.0.0.1 -P 3308 -u portal_user -p
```

When connecting from inside Docker network:
```bash
mysql -h mysql -P 3306 -u portal_user -p
```

## Changing Ports

To change ports, edit `docker-compose.yml`:

```yaml
services:
  frontend:
    ports:
      - "YOUR_PORT:80"  # Change YOUR_PORT to desired host port
  
  backend:
    ports:
      - "YOUR_PORT:8080"
  
  mysql:
    ports:
      - "YOUR_PORT:3306"
```

## Port Conflicts

If you encounter port conflicts, check which process is using the port:

```bash
# On macOS/Linux
lsof -i :8082
lsof -i :8080
lsof -i :3308

# On Windows
netstat -ano | findstr :8082
netstat -ano | findstr :8080
netstat -ano | findstr :3308
```

## Firewall Configuration

If deploying to production, ensure these ports are:
- **8082**: Open to public (frontend)
- **8080**: Protected, only accessible from frontend or trusted sources
- **3308**: Blocked from public access, only accessible from backend
