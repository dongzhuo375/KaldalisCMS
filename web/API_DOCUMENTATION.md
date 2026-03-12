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
POST /api/v1/admin/posts HTTP/1.1
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

The post API is split into two contracts:

1. **Public content API** under `/posts`
   - Read-only
   - Returns only `Published` content
2. **Admin content API** under `/admin/posts`
   - Requires authentication, authorization, and CSRF protection
   - `user` may create Draft posts
   - `admin` / `super_admin` can read drafts and perform publish/offline workflow transitions

### Public API

#### `GET /posts`

Returns the public post feed. This endpoint only includes posts whose `status` is `1` (`Published`).

**Responses:**

- `200 OK`: A list of published posts.
- `500 Internal Server Error`: Server error.

#### `GET /posts/:id`

Returns a single published post.

**URL Parameters:**

- `id` (integer, required): The ID of the post.

**Responses:**

- `200 OK`: The requested published post.
- `400 Bad Request`: Invalid post ID.
- `404 Not Found`: Post not found or not published.

### Admin API

All admin post routes are prefixed with `/api/v1/admin/posts` and require:

- a valid authenticated session
- a valid `X-CSRF-Token` header for state-changing requests
- route-specific role permissions

#### `GET /admin/posts`

Returns all posts, including drafts and published content, for management interfaces.

#### `GET /admin/posts/:id`

Returns a single post regardless of status for editing and moderation.

#### `POST /admin/posts`

Creates a new post.

**Permission:** authenticated `user`, `admin`, and `super_admin` roles may create drafts here.

**Important workflow rule:** newly created posts are always stored as `Draft`, even if the client tries to send another status.

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
- `403 Forbidden`: CSRF token invalid/missing or role lacks permission.
- `500 Internal Server Error`: Server error.

#### `PUT /admin/posts/:id`

Updates editable post content fields.

**Request Body:**

```json
{
  "title": "Updated Post Title",
  "content": "Updated content.",
  "cover": "/path/to/new-cover.jpg",
  "category_id": 2,
  "tags": [1, 4]
}
```

**Notes:**

- `status` is not changed by this endpoint.
- Publication state must be changed through the dedicated workflow endpoints below.

#### `POST /admin/posts/:id/publish`

Transitions a post from `Draft` to `Published`.

**Permission:** only `admin` and `super_admin` can publish.

**Responses:**

- `200 OK`: Post published successfully.
- `400 Bad Request`: Invalid transition or invalid post ID.
- `401 Unauthorized`: User not logged in.
- `403 Forbidden`: CSRF token invalid/missing or role lacks permission.
- `404 Not Found`: Post not found.

#### `POST /admin/posts/:id/draft`

Moves a post back to `Draft`.

This is the minimal "offline" action in the current workflow.

#### `DELETE /admin/posts/:id`

Deletes a post.

**Responses:**

- `204 No Content`: Post deleted successfully.
- `401 Unauthorized`: User not logged in.
- `403 Forbidden`: CSRF token invalid/missing or role lacks permission.
- `404 Not Found`: Post not found.