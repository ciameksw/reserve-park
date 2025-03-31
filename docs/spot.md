# Spot Service Documentation

## Endpoints

---

### 1. Create Spot
- **Method**: POST  
- **Endpoint**: `/spots`  
- **Description**: Creates a new parking spot.  
- **Request Body**:
    ```json
    {
        "latitude": 37.7749,
        "longitude": -122.4194,
        "price_per_hour": 5.5,
        "size": "medium",
        "type": "outdoor"
    }
    ```
- **Response**:
    - **201 Created**:
      123e4567-e89b-12d3-a456-426614174000
    - **400 Bad Request**: If the request body is invalid.
    - **500 Internal Server Error**: If there is an issue saving the spot.

---

### 2. Get Spot by ID
- **Method**: GET  
- **Endpoint**: `/spots/{id}`  
- **Description**: Retrieves a parking spot by its spot ID.  
- **Response**:
    - **200 OK**:
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
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the spot does not exist.
    - **500 Internal Server Error**: If there is an issue retrieving the spot.

---

### 3. Edit Spot
- **Method**: PATCH  
- **Endpoint**: `/spots`  
- **Description**: Updates an existing parking spot. Only the provided fields will be updated.  
- **Request Body**:
    ```json
    {
        "spot_id": "123e4567-e89b-12d3-a456-426614174000",
        "latitude": 37.7750,
        "price_per_hour": 6.0
    }
    ```
- **Response**:
    - **204 No Content**: If the update is successful.
    - **400 Bad Request**: If the request body is invalid.
    - **404 Not Found**: If the spot does not exist.
    - **500 Internal Server Error**: If there is an issue updating the spot.

---

### 4. Delete Spot
- **Method**: DELETE  
- **Endpoint**: `/spots/{id}`  
- **Description**: Deletes a parking spot by its spot ID.  
- **Response**:
    - **204 No Content**: If the deletion is successful.
    - **400 Bad Request**: If the `id` is missing or invalid.
    - **404 Not Found**: If the spot does not exist.
    - **500 Internal Server Error**: If there is an issue deleting the spot.

---

### 5. Get All Spots
- **Method**: GET  
- **Endpoint**: `/spots`  
- **Description**: Retrieves a list of all parking spots.  
- **Response**:
    - **200 OK**:
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
    - **500 Internal Server Error**: If there is an issue retrieving the spots.

---

### 6. Get Spot Price
- **Method**: GET  
- **Endpoint**: `/spots/price`  
- **Description**: Calculates the price for reserving a spot for a specific time range.  
- **Request Body**:
    ```json
    {
        "spot_id": "123e4567-e89b-12d3-a456-426614174000",
        "start_time": "2025-03-30T10:00:00Z",
        "end_time": "2025-03-30T12:00:00Z"
    }
    ```
- **Response**:
    - **200 OK**:
      ```json
      {
          "spot_id": "123e4567-e89b-12d3-a456-426614174000",
          "price": 11.0
      }
      ```
    - **400 Bad Request**: If the request body is invalid or the start time is after the end time.
    - **500 Internal Server Error**: If there is an issue calculating the price.

---

## MongoDB Document

### Spot Schema
```json
{
    "_id": "ObjectId",
    "spot_id": "string",
    "latitude": "float",
    "longitude": "float",
    "price_per_hour": "float",
    "size": "string", // e.g., "small", "medium", "large"
    "type": "string", // e.g., "indoor", "outdoor", "ev"
    "updated_at": "ISODate"
}
```