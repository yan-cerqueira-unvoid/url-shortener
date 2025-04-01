# URL Shortener Service

A simple and efficient URL shortening service built with Go and MongoDB.

## Features

- Shorten long URLs into compact, shareable links
- Custom short codes (optional)
- URL validation and normalization
- Automatic expiration of shortened URLs (default: 1 year)
- Click tracking for shortened URLs
- Configurable via environment variables

## Tech Stack

- Go with Gin web framework
- MongoDB for storage
- Docker-ready with configurable settings

## API Endpoints

- `GET /` - API information
- `POST /shorten` - Create a shortened URL
- `GET /:shortCode` - Redirect to the original URL

## Configuration

The service is configured via environment variables:

| Variable                | Description               | Default                   |
| ----------------------- | ------------------------- | ------------------------- |
| PORT                    | Server port               | 8080                      |
| MONGO_URI               | MongoDB connection string | mongodb://localhost:27017 |
| DB_NAME                 | Database name             | url_shortener             |
| URL_CODE_LENGTH         | Short code length         | 6                         |
| URL_DEFAULT_EXPIRY_DAYS | URL validity in days      | 365                       |
