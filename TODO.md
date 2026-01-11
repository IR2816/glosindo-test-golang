# TODO: Fix CORS Issues

## Information Gathered
- CORS middleware is configured in `middleware/cors.go` using gin-contrib/cors.
- Origins are set in `config/config.go` as hardcoded values: ["http://localhost:52302", "https://your-frontend-domain.com"].
- AllowCredentials is true, so origins cannot use "*".
- AllowHeaders include "Origin", "Content-Type", "Authorization".
- AllowMethods include GET, POST, PUT, DELETE, OPTIONS.
- Middleware is applied in `main.go`.

## Plan
- Update `config/config.go` to load CORSOrigins from environment variable `CORS_ORIGINS`, with default values.
- Update `middleware/cors.go` to include additional common headers like "Accept", "X-Requested-With", "X-CSRF-Token" for better compatibility.
- Ensure preflight OPTIONS requests are handled properly (already included in AllowMethods).

## Dependent Files to Edit
- `config/config.go`: Modify to load CORSOrigins from env.
- `middleware/cors.go`: Add more headers to AllowHeaders.

## Followup Steps
- Test the application to ensure CORS works for the specified origins.
- If needed, add more origins or adjust based on frontend requirements.
