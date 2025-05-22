# Facade Service Documentation

## Overview

The **facade service** acts as the main API gateway for the Reserve-Park system. It routes requests to the appropriate internal services (user, spot, reservation), handles authentication and authorization, and provides a unified interface for clients.

---

## Endpoints

---

### User Endpoints

#### Register User
- **POST** `/users/register`
- **Description:** Registers a new user in the system.
- **Request Body:**
    ```json
    {
      "username": "johndoe",
      "email": "johndoe@example.com",
      "password": "securepassword"
    }
    ```
- **Response:**
    - **201 Created**: User registered successfully.
    ```
    <USER_ID>
    ```
    - **400 Bad Request**: Invalid input.
    - **409 Conflict**: Username or email already exists.
    - **500 Internal Server Error**

---

#### Login
- **POST** `/users/login`
- **Description:** Authenticates a user and returns a JWT token.
- **Request Body:**
    ```json
    {
      "username": "johndoe",
      "password": "securepassword"
    }
    ```
- **Response:**
    - **200 OK**
      ```json
      {
        "jwt": "<JWT_TOKEN>"
      }
      ```
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Invalid credentials.
    - **500 Internal Server Error**

---

#### Get All Users (Admin)
- **GET** `/users`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns a list of all users (admin only).
- **Response:**
    - **200 OK**: List of users.
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
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Edit User
- **PATCH** `/users`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Updates user details (admin or self).
- **Request Body:**
    ```json
    {
      "user_id": "user-uuid",
      "username": "newusername",
      "email": "newemail@example.com",
      "password": "newpassword"
    }
    ```
- **Response:**
    - **204 No Content**: User updated.
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not allowed.
    - **404 Not Found**: If the user does not exist.
    - **409 Conflict**: If the username or email already exists.
    - **500 Internal Server Error**

---

#### Edit User Role (Admin)
- **PATCH** `/users/role`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Changes a user's role (admin only).
- **Request Body:**
    ```json
    {
      "user_id": "user-uuid",
      "role": "admin"
    }
    ```
- **Response:**
    - **204 No Content**: Role updated.
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not allowed.
    - **500 Internal Server Error**

---

#### Get User by ID
- **GET** `/users/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns user details by ID (admin or self).
- **Response:**
    - **200 OK**: User details.
    ```json
    {
        "user_id": "123e4567-e89b-12d3-a456-426614174000",
        "username": "johndoe",
        "email": "johndoe@example.com",
        "role": "user",
        "updated_at": "2025-03-30T10:00:00Z"
    }
    ```
    - **401 Unauthorized**: Not allowed.
    - **404 Not Found**: If the user does not exist.
    - **500 Internal Server Error**

---

#### Delete User by ID
- **DELETE** `/users/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Deletes a user by ID (admin or self).
- **Response:**
    - **204 No Content**: User deleted.
    - **401 Unauthorized**: Not allowed.
    - **404 Not Found**: If the user does not exist.
    - **500 Internal Server Error**

---

### Spot Endpoints

#### Get Available Spots
- **GET** `/spots/available`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns available spots for a given time range and spot IDs.
- **Request Body:**
    ```json
    {
      "spot_ids": ["spot1", "spot2"],
      "start_time": "2025-05-22T10:00:00Z",
      "end_time": "2025-05-22T12:00:00Z"
    }
    ```
- **Response:**
    - **200 OK**: List of available spots.
    ```
    [
    "spot1"
    ]
    ```
    - **400 Bad Request**: Some spots do not exist or invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Get Spot Price
- **GET** `/spots/price`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Calculates the price for a spot and time range.
- **Request Body:**
    ```json
    {
      "spot_id": "spot1",
      "start_time": "2025-05-22T10:00:00Z",
      "end_time": "2025-05-22T12:00:00Z"
    }
    ```
- **Response:**
    - **200 OK**
      ```json
      {
        "spot_id": "spot1",
        "price": 15.0
      }
      ```
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Get All Spots
- **GET** `/spots`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Retrieves a list of all parking spots.
- **Response:**
    - **200 OK**: List of spots.
    ```json
    [
        {
            "spot_id": "123e4567-e89b-12d3-a456-426614174000",
            "latitude": 37.7749,
            "longitude": -122.4194,
            "price_per_hour": 5.5,
            "size": "medium",
            "type": "outdoor",
            "updated_at": "2025-03-30T10:00:00Z"
        },
        {
            "spot_id": "456e7890-e12b-34d5-a678-426614174001",
            "latitude": 37.7750,
            "longitude": -122.4195,
            "price_per_hour": 6.0,
            "size": "large",
            "type": "indoor",
            "updated_at": "2025-03-30T11:00:00Z"
        }
    ]
    ```
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Get Spot by ID
- **GET** `/spots/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Retrieves information about a specific parking spot.
- **Response:**
    - **200 OK**: Spot info.
    ```json
    {
        "spot_id": "123e4567-e89b-12d3-a456-426614174000",
        "latitude": 37.7749,
        "longitude": -122.4194,
        "price_per_hour": 5.5,
        "size": "medium",
        "type": "outdoor",
        "updated_at": "2025-03-30T10:00:00Z"
    }
    ```
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Spot does not exist.
    - **500 Internal Server Error**

---

#### Delete Spot by ID (Admin)
- **DELETE** `/spots/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Deletes a parking spot by ID (admin only).
- **Response:**
    - **204 No Content**: Spot deleted.
    - **400 Bad Request**: `id` is missing or invalid.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Spot does not exist.
    - **500 Internal Server Error**

---

#### Add Spot (Admin)
- **POST** `/spots`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Adds a new parking spot (admin only).
- **Request Body:**
    ```json
    {
      "latitude": 42.712716,
      "longitude": -74.005974,
      "price_per_hour": 15.50,
      "size": "medium",
      "type": "outdoor"
    }
    ```
- **Response:**
    - **201 Created**: Spot created.
    ```
    123e4567-e89b-12d3-a456-426614174000
    ```
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Edit Spot (Admin)
- **PATCH** `/spots`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Edits an existing parking spot (admin only).
- **Request Body:**
    ```json
    {
      "spot_id": "spot1",
      "size": "large"
    }
    ```
- **Response:**
    - **204 No Content**: Spot updated.
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Spot does not exist.
    - **500 Internal Server Error**

---

### Reservation Endpoints

#### Get All Reservations (Admin)
- **GET** `/reservations`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns all reservations (admin only).
- **Response:**
    - **200 OK**: List of reservations.
    ```json
    [
        {
            "reservation_id": "123e4567-e89b-12d3-a456-426614174000",
            "user_id": "user123",
            "spot_id": "spot456",
            "start_time": "2025-03-30T10:00:00Z",
            "end_time": "2025-03-30T12:00:00Z",
            "status": "valid",
            "price_paid": 50.0,
            "updated_at": "2025-03-30T10:00:00Z"
        },
        {
            "reservation_id": "456e7890-e12b-34d5-a678-426614174001",
            "user_id": "user456",
            "spot_id": "spot789",
            "start_time": "2025-03-31T14:00:00Z",
            "end_time": "2025-03-31T16:00:00Z",
            "status": "valid",
            "price_paid": 60.0,
            "updated_at": "2025-03-31T14:00:00Z"
        }
    ]
    ```
    - **401 Unauthorized**: Not authenticated.
    - **500 Internal Server Error**

---

#### Get Reservations by Spot (Admin)
- **GET** `/reservations/spot/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns all reservations for a specific spot (admin only).
- **Response:**
    - **200 OK**: List of reservations.
    ```json
    [
        {
            "reservation_id": "123e4567-e89b-12d3-a456-426614174000",
            "user_id": "user123",
            "spot_id": "spot456",
            "start_time": "2025-03-30T10:00:00Z",
            "end_time": "2025-03-30T12:00:00Z",
            "status": "valid",
            "price_paid": 50.0,
            "updated_at": "2025-03-30T10:00:00Z"
        }
    ]
    ```
    - **400 Bad Request**: `id` is missing or invalid.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: No reservations exist for the spot.
    - **500 Internal Server Error**

---

#### Delete Reservation by ID (Admin)
- **DELETE** `/reservations/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Deletes a reservation by ID (admin only).
- **Response:**
    - **204 No Content**: Reservation deleted.
    - **400 Bad Request**: `id` is missing or invalid.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Reservation does not exist.
    - **500 Internal Server Error**

---

#### Get Reservations by User ID
- **GET** `/reservations/user/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns all reservations for a specific user (admin or self).
- **Response:**
    - **200 OK**: List of reservations.
    ```json
    [
        {
            "reservation_id": "123e4567-e89b-12d3-a456-426614174000",
            "user_id": "user123",
            "spot_id": "spot456",
            "start_time": "2025-03-30T10:00:00Z",
            "end_time": "2025-03-30T12:00:00Z",
            "status": "valid",
            "price_paid": 50.0,
            "updated_at": "2025-03-30T10:00:00Z"
        }
    ]
    ```
    - **400 Bad Request**: `id` is missing or invalid.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: No reservations exist for the user.
    - **500 Internal Server Error**

---

#### Get Reservation by ID
- **GET** `/reservations/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Returns reservation details by ID (admin or self).
- **Response:**
    - **200 OK**: Reservation details.
    ```json
    {
        "reservation_id": "123e4567-e89b-12d3-a456-426614174000",
        "user_id": "user123",
        "spot_id": "spot456",
        "start_time": "2025-03-30T10:00:00Z",
        "end_time": "2025-03-30T12:00:00Z",
        "status": "valid",
        "price_paid": 50.0,
        "updated_at": "2025-03-30T10:00:00Z"
    }
    ```
    - **400 Bad Request**: `id` is missing or invalid.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Reservation does not exist.
    - **500 Internal Server Error**

---

#### Add Reservation
- **POST** `/reservations`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Creates a new reservation (admin or self).
- **Request Body:**
    ```json
    {
      "user_id": "user-uuid",
      "spot_id": "spot-uuid",
      "start_time": "2025-05-22T10:00:00Z",
      "end_time": "2025-05-22T12:00:00Z",
      "price_paid": 35.95
    }
    ```
- **Response:**
    - **201 Created**: Reservation created.
    ```
    ff360c0a-6502-46bf-a8be-60807f142ab8
    ```
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Spot does not exist.
    - **409 Conflict**: Spot is not available in the provided timeframe.
    - **500 Internal Server Error**

---

#### Edit Reservation
- **PATCH** `/reservations`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Edits an existing reservation (admin or self).
- **Request Body:**
    ```json
    {
      "reservation_id": "reservation-uuid",
      "user_id": "user-uuid",
      "spot_id": "spot-uuid",
      "start_time": "2025-05-22T11:00:00Z",
      "end_time": "2025-05-22T13:00:00Z",
      "price_paid": 7.95
    }
    ```
- **Response:**
    - **204 No Content**: Reservation updated.
    - **400 Bad Request**: Invalid input.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Reservation or spot does not exist.
    - **409 Conflict**: Spot is not available in the updated timeframe.
    - **500 Internal Server Error**

---

#### Cancel Reservation
- **PATCH** `/reservations/cancel/{id}`
- **Headers:**  
    - `Authorization: Bearer <JWT_TOKEN>`
- **Description:** Cancels a reservation by ID (admin or self).
- **Response:**
    - **204 No Content**: Reservation canceled.
    - **401 Unauthorized**: Not authenticated.
    - **404 Not Found**: Reservation does not exist.
    - **409 Conflict**: Reservation already canceled.
    - **500 Internal Server Error**

---

## Notes

- All endpoints that require authentication expect a JWT token in the `Authorization` header.
- The facade service handles routing, validation, and authorization for all requests.