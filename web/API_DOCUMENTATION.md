# KaldalisCMS API Documentation

This document provides a detailed description of the KaldalisCMS API. It is intended for front-end developers who need to interact with the CMS backend.

## Base URL

All API endpoints are prefixed with `/api/v1`.

**Example:** `http://yourdomain.com/api/v1`

## Authentication

The API uses **HttpOnly Cookies** for secure authentication and a **Double Submit Cookie** pattern for CSRF protection.

### How to Authenticate

1.  Call the `POST /users/login` endpoint.
2.  On success, the server will set the following cookies:
    *   **Auth Cookie** (`access_token`): Contains the signed session token (HttpOnly).
    *   **CSRF Cookie** (`csrf_token`): Contains the CSRF token (Readable by JavaScript).
    *   **Role Cookie** (`user_role`): Contains the user's role (Readable by JavaScript).
3.  Browsers will automatically send these cookies with subsequent requests.

### CSRF Protection

For all state-changing requests (POST, PUT, DELETE), you **must** include the CSRF token in the request headers.

1.  Read the value from the `csrf_token` cookie.
2.  Add it to the request header: `X-CSRF-Token: <value>`.

**Example:**
```http
POST /api/v1/posts HTTP/1.1
Cookie: access_token=...; csrf_token=abc123...
X-CSRF-Token: abc123...
```

---

## Users API

The Users API provides endpoints for user registration and login.

### `POST /users/register`

Registers a new user.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com"
}
```

**Responses:**

- `201 Created`: User created successfully.
  ```json
  {
    "message": "User created successfully"
  }
  ```
- `400 Bad Request`: Invalid request body.
- `500 Internal Server Error`: Server error.

### `POST /users/login`

Logs in a user and establishes a session via cookies.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Responses:**

- `200 OK`: Login successful. Sets `Set-Cookie` headers.
  ```json
  {
    "message": "Login successful",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user"
    },
    "expires_at": "2023-10-28T10:00:00Z"
  }
  ```
- `401 Unauthorized`: Invalid username or password.
- `400 Bad Request`: Invalid request body.
- `500 Internal Server Error`: Server configuration error.

### `POST /users/logout`

Logs out the user and clears the session cookies. Requires authentication.

**Responses:**

- `200 OK`: Logout successful.
  ```json
  {
    "message": "logged out"
  }
  ```

---

## Posts API

The Posts API provides endpoints for managing blog posts.

### `GET /posts`

Retrieves a list of all posts. Public endpoint.

**Responses:**

- `200 OK`: A list of posts.
  ```json
  [
    {
      "id": 1,
      "title": "My First Post",
      "slug": "my-first-post",
      "content": "This is the content of my first post.",
      "cover": "/path/to/cover.jpg",
      "status": 1,
      "author": {
        "id": 1,
        "username": "author_name"
      },
      "category": {
        "id": 1,
        "name": "Tech"
      },
      "tags": [
        {
          "id": 1,
          "name": "Golang"
        }
      ],
      "created_at": "2023-10-27T10:00:00Z",
      "updated_at": "2023-10-27T10:00:00Z"
    }
  ]
  ```
- `500 Internal Server Error`: Server error.

### `GET /posts/:id`

Retrieves a single post by its ID. Public endpoint.

**URL Parameters:**

- `id` (integer, required): The ID of the post.

**Responses:**

- `200 OK`: The requested post.
  ```json
  {
    "id": 1,
    "title": "My First Post",
    "slug": "my-first-post",
    "content": "This is the content of my first post.",
    "cover": "/path/to/cover.jpg",
    "status": 1,
    "author": {
      "id": 1,
      "username": "author_name"
    },
    "category": {
      "id": 1,
      "name": "Tech"
    },
    "tags": [
      {
        "id": 1,
        "name": "Golang"
      }
    ],
    "created_at": "2023-10-27T10:00:00Z",
    "updated_at": "2023-10-27T10:00:00Z"
  }
  ```
- `400 Bad Request`: Invalid post ID.
- `404 Not Found`: Post not found.

### `POST /posts`

Creates a new post. (Authentication required, CSRF token required)

**Request Body:**

```json
{
  "title": "New Post Title",
  "content": "Content of the new post.",
  "cover": "/path/to/image.jpg",
  "category_id": 2,
  "tags": [1, 3]
}
```

**Responses:**

- `201 Created`: Post created successfully.
- `400 Bad Request`: Invalid request body.
- `401 Unauthorized`: User not logged in.
- `403 Forbidden`: CSRF token invalid or missing.
- `500 Internal Server Error`: Server error.

### `PUT /posts/:id`

Updates an existing post. (Authentication required, CSRF token required)

**URL Parameters:**

- `id` (integer, required): The ID of the post to update.

**Request Body:**

```json
{
  "title": "Updated Post Title",
  "content": "Updated content.",
  "cover": "/path/to/new-cover.jpg",
  "category_id": 2,
  "tags": [1, 4],
  "status": 1
}
```

**Responses:**

- `200 OK`: Post updated successfully.
- `400 Bad Request`: Invalid post ID or request body.
- `401 Unauthorized`: User not logged in.
- `403 Forbidden`: CSRF token invalid or missing.
- `404 Not Found`: Post not found.

### `DELETE /posts/:id`

Deletes a post. (Authentication required, CSRF token required)

**URL Parameters:**

- `id` (integer, required): The ID of the post to delete.

**Responses:**

- `204 No Content`: Post deleted successfully.
- `400 Bad Request`: Invalid post ID.
- `401 Unauthorized`: User not logged in.
- `403 Forbidden`: CSRF token invalid or missing.
- `404 Not Found`: Post not found.