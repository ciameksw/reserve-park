# Reservation Service Documentation

## Endpoints

---

### 1. Create Reservation
- **Method**: POST  
- **Endpoint**: `/reservations`  
- **Description**: Creates a new reservation.  
- **Request Body**:
    ```json
    {
        "user_id": "user123",
        "spot_id": "spot456",
        "start_time": "2025-03-30T10:00:00Z",
        "end_time": "2025-03-30T12:00:00Z",
        "status": "valid",
        "price_paid": 50.0
    }
    ```
- **Response**:
    - **201 Created**:
      ```json
      {
          "reservation_id": "123e4567-e89b-12d3-a456-426614174000"
      }
      ```
    - **400 Bad Request**: If the request body is invalid.
    - **409 Conflict**: If the spot is not available in the provided timeframe.
    - **500 Internal Server Error**: If there is an issue saving the reservation.

---

### 2. Get Reservation by ID
- **Method**: GET  
- **Endpoint**: `/reservations/{id}`  
- **Description**: Retrieves a reservation by its reservation ID.  
- **Response**:
    - **200 OK**:
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
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the reservation does not exist.
    - **500 Internal Server Error**: If there is an issue retrieving the reservation.

---

### 3. Edit Reservation
- **Method**: PATCH  
- **Endpoint**: `/reservations`  
- **Description**: Updates an existing reservation. Only the provided fields will be updated.  
- **Request Body**:
    ```json
    {
        "reservation_id": "123e4567-e89b-12d3-a456-426614174000",
        "start_time": "2025-03-30T11:00:00Z",
        "end_time": "2025-03-30T13:00:00Z"
    }
    ```
- **Response**:
    - **204 No Content**: If the update is successful.
    - **400 Bad Request**: If the request body is invalid.
    - **404 Not Found**: If the reservation does not exist.
    - **409 Conflict**: If the spot is not available in the updated timeframe.
    - **500 Internal Server Error**: If there is an issue updating the reservation.

---

### 4. Delete Reservation
- **Method**: DELETE  
- **Endpoint**: `/reservations/{id}`  
- **Description**: Deletes a reservation by its reservation ID.  
- **Response**:
    - **204 No Content**: If the deletion is successful.
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the reservation does not exist.
    - **500 Internal Server Error**: If there is an issue deleting the reservation.

---

### 5. Get All Reservations
- **Method**: GET  
- **Endpoint**: `/reservations`  
- **Description**: Retrieves a list of all reservations.  
- **Response**:
    - **200 OK**:
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
    - **500 Internal Server Error**: If there is an issue retrieving the reservations.

---

### 6. Get Reservations by User ID
- **Method**: GET  
- **Endpoint**: `/reservations/user/{id}`  
- **Description**: Retrieves all reservations for a specific user.  
- **Response**:
    - **200 OK**:
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
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If no reservations exist for the user.
    - **500 Internal Server Error**: If there is an issue retrieving the reservations.

---

### 7. Get Reservations by Spot ID
- **Method**: GET  
- **Endpoint**: `/reservations/spot/{id}`  
- **Description**: Retrieves all reservations for a specific spot.  
- **Response**:
    - **200 OK**:
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
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If no reservations exist for the spot.
    - **500 Internal Server Error**: If there is an issue retrieving the reservations.

---

### 8. Check Availability
- **Method**: GET  
- **Endpoint**: `/reservations/availability/check`  
- **Description**: Checks the availability of spots for a specific time range.  
- **Request Body**:
    ```json
    {
        "spot_ids": ["spot456", "spot789"],
        "start_time": "2025-03-30T10:00:00Z",
        "end_time": "2025-03-30T12:00:00Z"
    }
    ```
- **Response**:
    - **200 OK**:
      ```json
      ["spot456"]
      ```
    - **400 Bad Request**: If the request body is invalid.
    - **500 Internal Server Error**: If there is an issue checking availability.

---

## MongoDB Document

### Reservation Schema
```json
{
    "_id": "ObjectId",
    "reservation_id": "string",
    "user_id": "string",
    "spot_id": "string",
    "start_time": "ISODate",
    "end_time": "ISODate",
    "status": "string", // e.g., "valid", "canceled"
    "price_paid": "float",
    "updated_at": "ISODate"
}
```