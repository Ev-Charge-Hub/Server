# **EV-Charge-Hub Backend API Documentation**

## 📋 **Project Overview**

This project is a backend system for managing EV charging stations. It allows users to:

* Register and log in
* Retrieve details of all available EV stations
* Filter stations by company, type, and search keyword
* Retrieve a station by its unique ID

## **🏛️ Architecture: Clean Architecture**

This project follows **Clean Architecture** principles, which promote separation of concerns, maintainability, testability, and scalability of the system.

### **Clean Architecture Layers**
![image](https://github.com/user-attachments/assets/9ba3d515-64e1-4fa2-aaf2-b443a9c526f3)

1. **Delivery Layer (Controllers/Handlers):** Handles incoming HTTP requests, processes input, and returns responses.
2. **Use Case Layer (Business Logic):** Contains application-specific business rules and orchestrates the interaction between repository and response.
3. **Repository Layer (Data Access):** Interacts with external systems like databases or external APIs to store and retrieve data.
4. **Domain Layer (Entities/Models):** Represents core business entities and models with minimal dependencies.

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

## 📚 **API Endpoints**

### **1. User Management**

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/users/register` | Register a new user |
| POST | `/users/login` | Login with username/email |

#### 📋 **Register User**

* **URL:** `POST /users/register`
* **Body:**

  ```
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123",
    "role": "USER"
  }
  ```
* **Response:**

  ```
  {
    "message": "User registered successfully"
  }
  ```

#### 📋 **Login**

* **URL:** `POST /users/login`
* **Body:**

  ```
  {
    "username_or_email": "john@example.com",
    "password": "secret123"
  }
  ```
* **Response:**

  ```
  {
    "token": "your_jwt_token"
  }
  ```

---

### **2. EV Station Management**

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/stations` | Get all stations |
| GET | `/stations/filter` | Filter stations by parameters |
| GET | `/stations/:id` | Get station by ID |

#### 📋 **Get All Stations**

* **URL:** `GET /stations`
* **Response:**

  ```
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
  * `company`: Filter by company name (option)
  * `type`: Filter by connector type (option)
  * `search`: Search by station name (option)
  * `plug_name`: Filter by plug name (optional)
  * `status`: Filter by `open` / `closed` status of stations (optional) 
* Example
  * `GET /stations/filter?company=Caltex EV&type=AC&search=สถานีชาร์จรถยนต์ไฟฟ้า เซ็นทรัลเวิลด์`


* **Response:**

  ```
  [
    {
      "id": "3124AKCac,c9d6b1",
      "station_id": "ST001",
      "name": "สถานีชาร์จรถยนต์ไฟฟ้า เซ็นทรัลเวิลด์",
      "latitude": 13.7563,
      "longitude": 100.5018,
      "company": "Caltex EV",
      "status": {
        "open_hours": "08:00",
        "close_hours": "20:00",
        "is_open": true
      },
      "connectors": [
        {
          "connector_id": "C001",
          "type": "AC",
          "price_per_unit": 3.5,
          "power_output": 22.0,
          "is_available": true
        }
      ]
    },
    {
      "id": "67a0f307b60606077d9dc993",
      "station_id": "ST069",
      "name": "สถานีชาร์จรถยนต์ไฟฟ้า ห้างฉัตร ลำปาง",
      "latitude": 18.298543,
      "longitude": 99.30619,
      "company": "SUSCO EV",
      "status": {
          "open_hours": "06:00",
          "close_hours": "23:00",
          "is_open": true
      },
      "connectors": [
          {
              "connector_id": "CT0690",
              "type": "DC",
              "price_per_unit": 6.6,
              "power_output": 100,
              "is_available": true
          },
          {
              "connector_id": "CT0691",
              "type": "AC",
              "price_per_unit": 4.3,
              "power_output": 22,
              "is_available": false
          },
          {
              "connector_id": "CT0692",
              "type": "DC",
              "price_per_unit": 7.2,
              "power_output": 150,
              "is_available": true
          }
      ]
  }
  ]
  ```

#### 

#### 📋 **Get Station by ID**

* **URL:** `GET /stations/:id`
* **Path Parameter:**
  * `id`: The MongoDB ObjectID of the station
* Example
  * `GET /stations/63f5a01c8f7e3f65b4c9d6b1`


* **Response:**

```
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
```

#### 

#### 📋 **SetBooking**

* **URL:** `PUT /stations/set-booking`
* **Body:**
```
  {
    "connector_id": "CT0010",
    "username": "note",
    "booking_end_time": "2025-04-20T15:00:00"
}
```
* Example
  * `GET /stations/63f5a01c8f7e3f65b4c9d6b1`


* **Response:**

```
{
    "message": "Booking successfully added"
}
```




---

## 🔐 **JWT Authentication**

The system uses **JWT (JSON Web Token)** for user authentication. When logging in, the server will return a JWT token that must be included in the **Authorization** header of requests that require authentication.

Example:

`Authorization: Bearer <your_jwt_token>`

---

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
