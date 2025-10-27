# API Documentation

## Overview

This document provides detailed information about the Cutter Project API endpoints, request/response formats, authentication, and error handling.

## Base URL

```
http://localhost:8080/api
```

## Authentication

The API uses JWT (JSON Web Token) for authentication. After logging in, you'll receive an access token that must be included in the `Authorization` header of all protected requests.

```
Authorization: Bearer <access_token>
```

## Response Format

All API responses follow a consistent JSON format:

### Success Response

```json
{
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "error": "Error message",
  "details": "Additional error details (if available)",
  "request_id": "Unique request identifier (if available)"
}
```

## Status Codes

The API uses standard HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service unavailable

## Endpoints

### Authentication

#### Register User

Register a new user account.

- **URL**: `/users/register`
- **Method**: `POST`
- **Auth Required**: No

**Request Body**:

```json
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response**:

```json
{
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

**Error Codes**:

- `400 Bad Request` - Invalid input data
- `409 Conflict` - User with this email already exists

#### Login User

Authenticate a user and receive an access token.

- **URL**: `/users/login`
- **Method**: `POST`
- **Auth Required**: No

**Request Body**:

```json
{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}
```

**Response**:

```json
{
  "data": {
    "access_token": "jwt_access_token",
    "refresh_token": "jwt_refresh_token",
    "expires_in": 3600,
    "token_type": "Bearer"
  }
}
```

**Error Codes**:

- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Invalid credentials

### User Management

#### Get User

Retrieve user information by ID.

- **URL**: `/user/:id`
- **Method**: `GET`
- **Auth Required**: Yes

**Parameters**:

- `id` (path, required) - User ID

**Response**:

```json
{
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

**Error Codes**:

- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

#### Update User

Update user information.

- **URL**: `/user/:id`
- **Method**: `PUT`
- **Auth Required**: Yes

**Parameters**:

- `id` (path, required) - User ID

**Request Body**:

```json
{
  "name": "Jane Doe"
}
```

**Response**:

```json
{
  "data": {
    "id": "uuid",
    "name": "Jane Doe",
    "email": "john.doe@example.com",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

**Error Codes**:

- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

#### Change Password

Change user password.

- **URL**: `/user/:id/password`
- **Method**: `PUT`
- **Auth Required**: Yes

**Parameters**:

- `id` (path, required) - User ID

**Request Body**:

```json
{
  "current_password": "securepassword123",
  "new_password": "newsecurepassword456"
}
```

**Response**:

```json
{
  "data": {
    "message": "Password changed successfully"
  }
}
```

**Error Codes**:

- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

#### Delete User

Delete a user account.

- **URL**: `/user/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes

**Parameters**:

- `id` (path, required) - User ID

**Response**:

```json
{
  "data": {
    "message": "User deleted successfully"
  }
}
```

**Error Codes**:

- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - User not found

### Health Check

#### Health Check

Check the health status of the API and its dependencies.

- **URL**: `/health`
- **Method**: `GET`
- **Auth Required**: No

**Response**:

```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2023-01-01T00:00:00Z",
    "version": "1.0.0",
    "database": "healthy",
    "redis": "healthy"
  }
}
```

**Error Codes**:

- `503 Service Unavailable` - Service or dependencies unhealthy

## Rate Limiting

The API implements rate limiting to prevent abuse:

- Public endpoints: 60 requests per minute
- Protected endpoints: 120 requests per minute

Rate limit headers are included in responses:

- `X-RateLimit-Limit` - Maximum number of requests allowed
- `X-RateLimit-Remaining` - Remaining requests in the current window
- `X-RateLimit-Reset` - Time when the current window resets (UTC timestamp)

## Error Handling

The API provides detailed error information:

```json
{
  "error": "Error message",
  "details": "Additional error details (if available)",
  "request_id": "Unique request identifier (if available)"
}
```

Common error types:

- `Bad Request` - Invalid input data
- `Unauthorized` - Authentication required or invalid
- `Forbidden` - Insufficient permissions
- `Not Found` - Resource not found
- `Conflict` - Resource conflict (e.g., duplicate email)
- `Unprocessable Entity` - Validation error
- `Too Many Requests` - Rate limit exceeded
- `Internal Server Error` - Server error
- `Service Unavailable` - Service or dependencies unhealthy

## Examples

### Register a User

```bash
curl -X POST http://localhost:8080/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

### Login a User

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securepassword123"
  }'
```

### Get User Information

```bash
curl -X GET http://localhost:8080/api/user/123 \
  -H "Authorization: Bearer <access_token>"
```

### Update User Information

```bash
curl -X PUT http://localhost:8080/api/user/123 \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe"
  }'
```

### Change Password

```bash
curl -X PUT http://localhost:8080/api/user/123/password \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "current_password": "securepassword123",
    "new_password": "newsecurepassword456"
  }'
```

### Delete User

```bash
curl -X DELETE http://localhost:8080/api/user/123 \
  -H "Authorization: Bearer <access_token>"
```

### Health Check

```bash
curl -X GET http://localhost:8080/api/health