# KaldalisCMS API Tests

# NOTE:
# Post management routes now live under /api/v1/admin/posts.
# The login account used below must have admin or super_admin permissions.

# 1. Login as an admin user and capture auth tokens
POST http://localhost:8080/api/v1/users/login
Content-Type: application/json
{
  "username": "admin",
  "password": "password123"
}
HTTP 200
[Captures]
csrf_token: cookie "kaldalis_csrf"
auth_token: cookie "kaldalis_auth"

# 2. Create a new draft post through the admin API
POST http://localhost:8080/api/v1/admin/posts
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
Content-Type: application/json
{
  "title": "Minimal Post",
  "content": "This is a minimal test post."
}
HTTP 201

# 3. Get all admin posts and capture the ID of the first one
GET http://localhost:8080/api/v1/admin/posts
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
HTTP 200
[Asserts]
jsonpath "$" count > 0
[Captures]
new_post_id: jsonpath "$[0].id"

# 4. Verify the newly created post is still a draft in the admin detail view
GET http://localhost:8080/api/v1/admin/posts/{{new_post_id}}
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{new_post_id}}
jsonpath "$.status" == 0

# 5. Publish the draft so it becomes visible on the public API
POST http://localhost:8080/api/v1/admin/posts/{{new_post_id}}/publish
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
HTTP 200

# 6. Verify the public posts feed can now see the published post
GET http://localhost:8080/api/v1/posts
HTTP 200
[Asserts]
jsonpath "$[0].status" == 1

# 7. Update the published post through the admin API
PUT http://localhost:8080/api/v1/admin/posts/{{new_post_id}}
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
Content-Type: application/json
{
  "title": "Updated Minimal Post",
  "content": "Updated content content."
}
HTTP 200

# 8. Move the post back to draft (offline)
POST http://localhost:8080/api/v1/admin/posts/{{new_post_id}}/draft
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
HTTP 200

# 9. Delete the post from the admin API
DELETE http://localhost:8080/api/v1/admin/posts/{{new_post_id}}
Cookie: kaldalis_csrf={{csrf_token}}; kaldalis_auth={{auth_token}}
X-CSRF-Token: {{csrf_token}}
HTTP 204
