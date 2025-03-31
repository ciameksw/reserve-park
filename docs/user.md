# User Service Documentation

## Endpoints

---

### 1. Create User
- **Method**: POST  
- **Endpoint**: `/users`  
- **Description**: Creates a new user.  
- **Request Body**:
    ```json
    {
        "username": "johndoe",
        "email": "johndoe@example.com",
        "password": "securepassword",
        "role": "user"
    }
    ```
- **Response**:
    - **201 Created**:
      ```json
      {
          "user_id": "123e4567-e89b-12d3-a456-426614174000"
      }
      ```
    - **400 Bad Request**: If the request body is invalid.
    - **409 Conflict**: If the username or email already exists.
    - **500 Internal Server Error**: If there is an issue saving the user.

---

### 2. Get User by ID
- **Method**: GET  
- **Endpoint**: `/users/{id}`  
- **Description**: Retrieves a user by their user ID.  
- **Response**:
    - **200 OK**:
      ```json
      {
          "user_id": "123e4567-e89b-12d3-a456-426614174000",
          "username": "johndoe",
          "email": "johndoe@example.com",
          "role": "user",
          "updated_at": "2025-03-30T10:00:00Z"
      }
      ```
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the user does not exist.
    - **500 Internal Server Error**: If there is an issue retrieving the user.

---

### 3. Edit User
- **Method**: PATCH  
- **Endpoint**: `/users`  
- **Description**: Updates an existing user. Only the provided fields will be updated.  
- **Request Body**:
    ```json
    {
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "new_username",
        "email": "new_email@example.com",
        "password": "new_securepassword",
        "role": "admin"
    }
    ```
- **Response**:
    - **204 No Content**: If the update is successful.
    - **400 Bad Request**: If the request body is invalid.
    - **404 Not Found**: If the user does not exist.
    - **409 Conflict**: If the username or email already exists.
    - **500 Internal Server Error**: If there is an issue updating the user.

---

### 4. Delete User
- **Method**: DELETE  
- **Endpoint**: `/users/{id}`  
- **Description**: Deletes a user by their user ID.  
- **Response**:
    - **204 No Content**: If the deletion is successful.
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the user does not exist.
    - **500 Internal Server Error**: If there is an issue deleting the user.

---

### 5. Get All Users
- **Method**: GET  
- **Endpoint**: `/users`  
- **Description**: Retrieves a list of all users.  
- **Response**:
    - **200 OK**:
      ```json
      [
          {
              "user_id": "123e4567-e89b-12d3-a456-426614174000",
              "username": "johndoe",
              "email": "johndoe@example.com",
              "role": "user",
              "updated_at": "2025-03-30T10:00:00Z"
          },
          {
              "user_id": "456e7890-e12b-34d5-a678-426614174001",
              "username": "janedoe",
              "email": "janedoe@example.com",
              "role": "admin",
              "updated_at": "2025-03-30T11:00:00Z"
          }
      ]
      ```
    - **500 Internal Server Error**: If there is an issue retrieving the users.

---

### 6. Login
- **Method**: POST  
- **Endpoint**: `/users/login`  
- **Description**: Authenticates a user and returns a JWT token.  
- **Request Body**:
    ```json
    {
        "username": "johndoe",
        "password": "securepassword"
    }
    ```
- **Response**:
    - **200 OK**:
      ```json
      {
          "jwt": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
      }
      ```
    - **400 Bad Request**: If the request body is invalid.
    - **401 Unauthorized**: If the username or password is incorrect.
    - **500 Internal Server Error**: If there is an issue generating the JWT token.

---

### 7. Authorize
- **Method**: GET  
- **Endpoint**: `/users/authorize`  
- **Description**: Validates a JWT token and returns the user's role.  
- **Headers**:
    - `Authorization`: `Bearer <JWT_TOKEN>`  
- **Response**:
    - **200 OK**:
      ```json
      {
          "role": "user"
      }
      ```
    - **401 Unauthorized**: If the token is missing, invalid, or expired.
    - **500 Internal Server Error**: If there is an issue validating the token.

---

## MongoDB Document

### User Schema
```json
{
    "_id": "ObjectId",
    "user_id": "string",
    "username": "string",
    "email": "string",
    "password_hash": "string",
    "role": "string", // e.g., "admin", "user"
    "updated_at": "ISODate"
}
```