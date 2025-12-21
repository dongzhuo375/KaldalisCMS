# KaldalisCMS API Tests

# 1. Register a new user
POST http://localhost:8080/api/v1/users/register
Content-Type: application/json
```json
{
 "username": "testuser",
 "password": "password123",
  "email": "test@example.com",
 "role": "user"
}
```tester
HTTP 201

# 2. Login as the new user and capture auth tokens
POST http://localhost:8080/api/v1/users/login
Content-Type: application/json
```json
{
  "username": "testuser",
  "password": "password123"
}
```
HTTP 200
[Captures]
csrf_token: cookie "kaldalis_csrf"
auth_token: cookie "kaldalis_auth"

# 3. Create a new blog post
POST http://localhost:8080/api/v1/posts
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
Content-Type: application/json
```json
{
  "title": "Minimal Post",
  "content": "This is a minimal test post."
}
```
HTTP 201

# 4. Get all posts and capture the ID of the first one
GET http://localhost:8080/api/v1/posts
HTTP 200
[Asserts]
jsonpath "$" count > 0
[Captures]
new_post_id: jsonpath "$[0].id"

# 5. Get the newly created post by its ID
GET http://localhost:8080/api/v1/posts/{{new_post_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{new_post_id}}
jsonpath "$.title" == "Minimal Post"

# 6. Update the created post
PUT http://localhost:8080/api/v1/posts/{{new_post_id}}
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
Content-Type: application/json
```json
{
  "title": "Updated Minimal Post"
}
```
HTTP 200

# 7. Verify the post was updated
GET http://localhost:8080/api/v1/posts/{{new_post_id}}
HTTP 200
[Asserts]
jsonpath "$.title" == "Updated Minimal Post"

# 8. Delete the post
DELETE http://localhost:8080/api/v1/posts/{{new_post_id}}
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
HTTP 204

# 9. Verify the post was deleted
GET http://localhost:8080/api/v1/posts/{{new_post_id}}
HTTP 404
