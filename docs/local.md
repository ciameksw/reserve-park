# Running Reserve-Park Locally

This guide explains how to start the Reserve Park system on your local machine, including restoring MongoDB data from a backup.

---

## 1. Build and Start All Services

Run the following command in the project root to build and start all services in detached mode:

```sh
docker compose up --build -d
```

This will start the MongoDB database and all microservices defined in `compose.yaml`.

## 2. Restore MongoDB Data from Backup

After the containers are running, restore the MongoDB databases using the backup provided in the `./backup` directory:

```sh
docker exec -it mongodb mongorestore /backup
```

## 3. Using the API

Once the services are running and the backup data is restored, you can interact with the system using the API endpoints provided by facade service.

**Default users available after restoring the backup:**

- **Admin user:**  
  - Username: `admin`  
  - Password: `admin`

- **Normal user:**  
  - Username: `user`  
  - Password: `user`

You can use these credentials to log in and test the API with different permission levels.

This command runs `mongorestore` inside the MongoDB container, importing all databases and collections from the backup.

## 4. Stopping the System

To stop all services, run:

```sh
docker compose down
```