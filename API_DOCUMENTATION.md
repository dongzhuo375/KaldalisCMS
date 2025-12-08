# KaldalisCMS API Documentation

This document provides a detailed description of the KaldalisCMS API. It is intended for front-end developers who need to interact with the CMS backend.

## Base URL

All API endpoints are prefixed with `/api/v1`.

**Example:** `http://yourdomain.com/api/v1`

## Authentication

Some endpoints require authentication using a JSON Web Token (JWT). To authenticate, you must first obtain a token by using the `POST /users/login` endpoint.

Once you have the token, you must include it in the `Authorization` header of your requests for protected endpoints.

**Example:** `Authorization: Bearer <your_jwt_token>`

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
  "email": "test@example.com",
  "role": "user"
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

Logs in a user and returns a JWT token.

**Request Body:**

```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Responses:**

- `200 OK`: Login successful.
  ```json
  {
    "message": "Login successful",
    "token": "<jwt_token>",
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com",
      "role": "user"
    }
  }
  ```
- `401 Unauthorized`: Invalid username or password.
- `400 Bad Request`: Invalid request body.

---

## Posts API

The Posts API provides endpoints for managing blog posts.

### `GET /posts`

Retrieves a list of all posts.

**Responses:**

- `200 OK`: A list of posts.
  ```json
  [
    {
      "ID": 1,
      "CreatedAt": "2023-10-27T10:00:00Z",
      "UpdatedAt": "2023-10-27T10:00:00Z",
      "Title": "My First Post",
      "Slug": "my-first-post",
      "Content": "This is the content of my first post.",
      "Cover": "/path/to/cover.jpg",
      "AuthorID": 1,
      "CategoryID": 1,
      "Status": 1
    }
  ]
  ```
- `500 Internal Server Error`: Server error.

### `GET /posts/:id`

Retrieves a single post by its ID.

**URL Parameters:**

- `id` (integer, required): The ID of the post.

**Responses:**

- `200 OK`: The requested post.
  ```json
  {
    "ID": 1,
    "CreatedAt": "2023-10-27T10:00:00Z",
    "UpdatedAt": "2023-10-27T10:00:00Z",
    "Title": "My First Post",
    "Slug": "my-first-post",
    "Content": "This is the content of my first post.",
    "Cover": "/path/to/cover.jpg",
    "AuthorID": 1,
    "CategoryID": 1,
    "Status": 1
  }
  ```
- `400 Bad Request`: Invalid post ID.
- `404 Not Found`: Post not found.

### `POST /posts`

Creates a new post. (Authentication required)

**Request Body:**

```json
{
  "Title": "New Post Title",
  "Slug": "new-post-title",
  "Content": "Content of the new post.",
  "Cover": "/path/to/image.jpg",
  "AuthorID": 1,
  "CategoryID": 2,
  "Status": 0
}
```

**Responses:**

- `201 Created`: Post created successfully.
- `400 Bad Request`: Invalid request body.
- `500 Internal Server Error`: Server error.

### `PUT /posts/:id`

Updates an existing post. (Authentication required)

**URL Parameters:**

- `id` (integer, required): The ID of the post to update.

**Request Body:**

```json
{
  "Title": "Updated Post Title",
  "Content": "Updated content."
}
```

**Responses:**

- `200 OK`: Post updated successfully.
- `400 Bad Request`: Invalid post ID or request body.
- `404 Not Found`: Post not found.

### `DELETE /posts/:id`

Deletes a post. (Authentication required)

**URL Parameters:**

- `id` (integer, required): The ID of the post to delete.

**Responses:**

- `204 No Content`: Post deleted successfully.
- `400 Bad Request`: Invalid post ID.
- `404 Not Found`: Post not found.
