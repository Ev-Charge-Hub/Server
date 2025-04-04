# **EV-Charge-Hub Backend API Documentation**

## 📋 Project Overview
EV-Charge-Hub is a backend system designed for managing electric vehicle (EV) charging stations. It enables user registration, station discovery, filtering by parameters and booking connectors. The system is built using the **Clean Architecture** design pattern for optimal maintainability and testability.

---

## 🧩 What is Clean Architecture?
Clean Architecture is a software design approach that separates the system into layers with clear boundaries and focuses on the **dependency inversion principle** — meaning that dependencies should always point inward toward the core business logic (domain). This architecture ensures that the core of the application is independent from frameworks, databases, and external tools.

Benefits include:
- Clear separation of concerns
- Easier testing and maintenance
- Flexibility in swapping infrastructure or delivery mechanisms

---
## 🧱 How It Works in This Project
The system consists of four key layers that match the Clean Architecture diagram shown below:

![image](https://github.com/user-attachments/assets/3b81a4cd-b5f4-4655-9c5a-8e4085d5b952)


### 🧠 1. Domain Layer (Core Business Logic)
This is the heart of the system, containing all critical business rules and enterprise logic.
- **Entity:** Models representing the business rules such as `Station`, `User`, `Booking`
- **Domain Model:** Defines core attributes and relationships of business objects
- 📁 Located at: `internal/domain/`

### 📤 2. Application Layer (UseCase Layer)
This layer orchestrates application-specific logic by invoking repositories and entities.
- **UseCase:** Implements business workflows and logic such as booking a connector
- Defines interfaces to be implemented by the repository layer
- 📁 Located at: `internal/usecase/`

### 📦 3. Data Layer (Infrastructure)
This layer contains implementations of external systems, especially databases.
- **Repository:** Implements the interfaces defined in the UseCase layer to access MongoDB
- 📁 Located at: `internal/repository/`

### 📬 4. Delivery Layer (Interface/Controllers)
The outermost layer that interacts with the client via HTTP, gRPC, etc.
- **Controller (Delivery):** Receives client requests, maps data to DTOs, and triggers UseCase logic
- 📁 Located at: `internal/delivery/`

---

### 🧾 Additional Components

- **DTO (Data Transfer Object):**
  - Defines how data is passed between the client and the server (request/response format)
  - 📁 Located at: `internal/dto/`

- **Middleware:**
  - Handles cross-cutting concerns like JWT authentication
  - 📁 Located at: `middleware/`

- **Utils:**
  - Helper functions for encryption, token generation, formatting, etc.
  - 📁 Located at: `utils/`

- **Routes:**
  - All route definitions and mapping to the appropriate handlers
  - 📁 Located at: `routes/`

- **Main Entry Point:**
  - Application bootstrap logic, initializing dependencies and servers
  - 📄 File: `main.go`

---

✅ This layered structure ensures the system is modular, testable, and easy to scale while following the principles of Clean Architecture.

## ⚙️ **Project Structure**

```
Project/
├── configs/                  # Configuration for DB connections
├── internal/
│   ├── delivery/             # HTTP Handlers (Controllers)
│   ├── repository/           # Data access logic (MongoDB)
│   ├── usecase/              # Business logic
│   ├── domain/               # Models for domain logic
│   └── dto/                  # DTOs (Data Transfer Objects)
├── middleware                # Middleware layer for verify before use restrict api
├── routes/                   # API routes
├── utils/                    # Utility functions (JWT, encryption)
└── main.go                   # Application entry point
```

---

## 🚀 **Getting Started**

### **1. Clone the repository**

- `git clone https://github.com/yourusername/ev-charge-hub.git`
- `cd ev-charge-hub`

### **2. Set up MongoDB**

Make sure MongoDB is installed and running on your local machine or connect to a remote instance.

Default connection string:

`mongodb://localhost:27017`

### **3. Configure environment variables**

Create a `.env` file in the project root and specify the following variables:

`MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=ev_charge_hub
JWT_SECRET=your_jwt_secret`

### **4. Install dependencies**

`go mod tidy`

### **5. Run the application**

`go run main.go`

Server will start at:

`http://localhost:8080`

---
## 🔐 JWT Authentication

All authenticated routes require a header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 📚 API Endpoints

### **1. User Management**

| Method | Endpoint            | Description              |
|--------|---------------------|--------------------------|
| POST   | `/users/register`   | Register a new user      |
| POST   | `/users/login`      | Login with username/email|

#### 📋 **Register User**
* **URL:** `POST /users/register`
* **Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secret123",
  "role": "USER"
}
```
* **Response:**
```json
{
  "message": "User registered successfully"
}
```

#### 📋 **Login**
* **URL:** `POST /users/login`
* **Body:**
```json
{
  "username_or_email": "john@example.com",
  "password": "secret123"
}
```
* **Response:**
```json
{
  "token": "your_jwt_token"
}
```

---

### **2. EV Station Management**

| Method | Endpoint              | Description             |
|--------|-----------------------|-------------------------|
| GET    | `/stations`           | Get all stations        |
| GET    | `/stations/filter`    | Filter stations         |
| GET    | `/stations/:id`       | Get station by ID       |
| POST   | `/stations/create`    | Create a new station    |
| PUT    | `/stations/:id`       | Update station info     |
| DELETE | `/stations/:id`       | Delete station          |

#### 📋 **Get All Stations**
* **URL:** `GET /stations`
* **Response:**
```json
[
  {
    "id": "63f5a01c8f7e3f65b4c9d6b1",
    "station_id": "ST001",
    "name": "EV Station Central Plaza",
    "latitude": 13.7563,
    "longitude": 100.5018,
    "company": "EV Company",
    "status": {
      "open_hours": "08:00",
      "close_hours": "20:00",
      "is_open": true
    },
    "connectors": [
        {
          "connector_id": "C001",
          "type": "Type2",
          "price_per_unit": 3.5,
          "power_output": 22.0,
          "is_available": true
        }
      ]
  }
]
```

#### 📋 **Filter Stations**
* **URL:** `GET /stations/filter`
* **Query Parameters:**
  - `company` (optional)
  - `type` (optional)
  - `search` (optional)
  - `plug_name` (optional)
  - `status` (`open` / `closed`, optional)
* **Response:**
```json
[
  {
    "id": "63f5a01c8f7e3f65b4c9d6b1",
    "station_id": "ST001",
    "name": "EV Station Central Plaza",
    "latitude": 13.7563,
    "longitude": 100.5018,
    "company": "EV Company",
    "status": {
      "open_hours": "08:00",
      "close_hours": "20:00",
      "is_open": true
    },
    "connectors": [
        {
          "connector_id": "C001",
          "type": "Type2",
          "price_per_unit": 3.5,
          "power_output": 22.0,
          "is_available": true
        }
      ]
  }
]
```

#### 📋 **Get Station by ID**
* **URL:** `GET /stations/:id`
* **Path Parameter:** `id`
* **Response:** station object 
```json
{
    "id": "67ee5a8fa3c75de5eb49699d",
    "station_id": "",
    "name": "Updated EV Station",
    "latitude": 13.75,
    "longitude": 100.5,
    "company": "EV Co Updated",
    "status": {
        "open_hours": "07:00",
        "close_hours": "22:00",
        "is_open": true
    },
    "connectors": [
        {
            "connector_id": "67ee5b77a3c75de5eb49699e",
            "type": "DC_FAST",
            "plug_name": "CCS",
            "price_per_unit": 8.5,
            "power_output": 120
        },
        {
            "connector_id": "67ee5b77a3c75de5eb49699f",
            "type": "AC_SLOW",
            "plug_name": "TYPE2",
            "price_per_unit": 5,
            "power_output": 22
        }
    ]
}
```

#### 📋 **Create Station**
* **URL:** `POST /stations/create`
* **Body:** full station object
```json
{
    "name": "test_company2_name4",
    "latitude": 13.304,
    "longitude": 100.46789,
    "company": "test_company2",
    "status": {
        "open_hours": "06:00",
        "close_hours": "23:00",
        "is_open": true
    },
    "connectors": [
        {
            "type": "DC",
            "plug_name": "CCS1",
            "price_per_unit": 7,
            "power_output": 100
        }
    ]
}

```
* **Response:** Created station details or success message

#### 📋 **Update Station**
* **URL:** `PUT /stations/:id`
* **Body:** full station object
```json
{
  "name": "Updated EV Station",
  "latitude": 13.7500,
  "longitude": 100.5000,
  "company": "EV Co Updated",
  "status": {
    "open_hours": "07:00",
    "close_hours": "22:00",
    "is_open": true
  },
  "connectors": [
    {
      "type": "DC_FAST",
      "plug_name": "CCS",
      "price_per_unit": 8.5,
      "power_output": 120
    },
    {
      "type": "AC_SLOW",
      "plug_name": "TYPE2",
      "price_per_unit": 5.0,
      "power_output": 22
    }
  ]
}
```

* **Response:** Updated info or success message
```json
{
    "message": "Station updated successfully",
    "station": {
        "id": "67ee5a8fa3c75de5eb49699d",
        "station_id": "",
        "name": "Updated EV Station",
        "latitude": 13.75,
        "longitude": 100.5,
        "company": "EV Co Updated",
        "status": {
            "open_hours": "07:00",
            "close_hours": "22:00",
            "is_open": true
        },
        "connectors": [
            {
                "connector_id": "67ee619c6cc170c4890c135e",
                "type": "DC_FAST",
                "plug_name": "CCS",
                "price_per_unit": 8.5,
                "power_output": 120
            },
            {
                "connector_id": "67ee619c6cc170c4890c135f",
                "type": "AC_SLOW",
                "plug_name": "TYPE2",
                "price_per_unit": 5,
                "power_output": 22
            }
        ]
    }
}
```

#### 📋 **Delete Station**
* **URL:** `DELETE /stations/:id`
* **Response:**
```json
{
  "message": "Station deleted successfully",
  "status": "success"
}
```

---

### **3. Booking Management**

| Method | Endpoint                         | Description                   |
|--------|----------------------------------|-------------------------------|
| PUT    | `/stations/set-booking`          | Set connector booking         |
| GET    | `/stations/booking/:username`    | Get booking by username       |
| GET    | `/stations/bookings/:username`   | Get all bookings for user     |

#### 📋 **Set Booking**
* **URL:** `PUT /stations/set-booking`
* **Body:** 
```json
{
  "connector_id": "CT0010",
  "username": "note",
  "booking_end_time": "2025-04-20T15:00:00"
}
```
* **Validation Rules Before Booking:**
	1. Reject if booking_end_time is in the past or now.
	2. Reject if user already has an active booking.
	3. Reject if connector is already booked by someone else.
	4. If all checks pass, create the booking.

* **Response:**
```json
{
  "message": "Booking successfully added"
}
```

#### 📋 **Get Booking by Username**
* **URL:** `GET /stations/booking/:username`
* **Response:** booking object
```json
{
  "username": "note",
  "booking_end_time": "2025-04-20T15:00:00"
}
```

#### 📋 **Get All Bookings by User**
* **URL:** `GET /stations/bookings/:username`
* **Response:** array of booking object
```json
[
  {
    "username": "note",
    "booking_end_time": "2025-04-20T15:00:00"
  },
  {
    "username": "note",
    "booking_end_time": "2025-05-20T15:00:00"
  }
]
```

---

### **4. Connector & User-specific Station API**

| Method | Endpoint                                | Description                          |
|--------|-----------------------------------------|--------------------------------------|
| GET    | `/stations/connector/:connector_id`     | Get station by connector ID          |
| GET    | `/stations/username/:username`          | Get stations by username             |

#### 📋 **Get Station by Connector ID**
* **URL:** `GET /stations/connector/:connector_id`
* **Response:** Full station object
```json
{
    "id": "67d7d957014efb03c444443a",
    "station_id": "ST001",
    "name": "สถานีชาร์จรถยนต์ไฟฟ้า เซ็นทรัลเวิลด์",
    "latitude": 13.746879,
    "longitude": 100.539742,
    "company": "PTT EV Station",
    "status": {
        "open_hours": "00:00",
        "close_hours": "23:59",
        "is_open": true
    },
    "connectors": [
        {
            "connector_id": "CT0010",
            "type": "DC",
            "plug_name": "CCS2",
            "price_per_unit": 6.5,
            "power_output": 150,
            "booking": {
                "username": "note",
                "booking_end_time": "2025-04-20T15:00:00"
            }
        },
        {
            "connector_id": "CT0011",
            "type": "AC",
            "plug_name": "Type 2",
            "price_per_unit": 4.5,
            "power_output": 22,
            "booking": {
                "username": "note",
                "booking_end_time": "2025-05-20T15:00:00"
            }
        },
        {
            "connector_id": "CT0012",
            "type": "DC",
            "plug_name": "CHAdeMO",
            "price_per_unit": 6.8,
            "power_output": 100,
            "booking": {
                "username": "MichaelBrown",
                "booking_end_time": "2025-03-17T10:05:08"
            }
        }
    ]
}
```

#### 📋 **Get Stations by Username**
* **URL:** `GET /stations/username/:username`
* **Response:** Full station object related to user
```json
{
    "id": "67d7d957014efb03c444443a",
    "station_id": "ST001",
    "name": "สถานีชาร์จรถยนต์ไฟฟ้า เซ็นทรัลเวิลด์",
    "latitude": 13.746879,
    "longitude": 100.539742,
    "company": "PTT EV Station",
    "status": {
        "open_hours": "00:00",
        "close_hours": "23:59",
        "is_open": true
    },
    "connectors": [
        {
            "connector_id": "CT0010",
            "type": "DC",
            "plug_name": "CCS2",
            "price_per_unit": 6.5,
            "power_output": 150,
            "booking": {
                "username": "note",
                "booking_end_time": "2025-04-20T15:00:00"
            }
        },
        {
            "connector_id": "CT0011",
            "type": "AC",
            "plug_name": "Type 2",
            "price_per_unit": 4.5,
            "power_output": 22,
            "booking": {
                "username": "note",
                "booking_end_time": "2025-05-20T15:00:00"
            }
        }
    ]
}

```


## 🛠 **Utilities**

* **Password encryption:** Uses `bcrypt` for hashing passwords before saving to the database.
* **JWT tokens:** Used for secure user authentication and session management.

## 🧾 **License**

This project is licensed under the MIT License.

---

## 🌟 **Acknowledgements**

Special thanks to the contributors and the Go community for making this project possible.

---

With this **README.md**, your documentation will guide users and developers to understand the project, configure it, and interact with it via the API.
