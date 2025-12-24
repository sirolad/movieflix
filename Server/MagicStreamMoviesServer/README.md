# MagicStreamMovies Server

MagicStreamMovies Server is a robust backend service built with **Go** and **Gin**, following modern clean architecture principles. It provides RESTful APIs for movie streaming features, including user authentication, movie management, genre filtering, and AI-powered recommendations.

## 🚀 Features

- **User Management**: Registration, Login (JWT), and Profile management.
- **Movie Management**: CRUD operations for movies.
- **Recommendations**: Personalized movie recommendations based on user favorites.
- **AI Integration**: Sentiment analysis and ranking for admin reviews using OpenAI.
- **Swagger Documentation**: Interactive API documentation.
- **Clean Architecture**: Separation of concerns (Handler -> Service -> Repository).

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin)
- **Database**: MongoDB
- **Authentication**: JWT (JSON Web Tokens)
- **Documentation**: Swagger (Swag)
- **Live Reload**: Air

## 📂 Project Structure

```
.
├── cmd
│   └── api
│       └── main.go           # Application entry point
├── internal
│   ├── config                # Configuration loader
│   ├── handler               # HTTP Handlers (Controllers)
│   ├── middleware            # HTTP Middleware (Auth, CORS)
│   ├── mocks                 # Mock implementations for testing
│   ├── models                # Data structures
│   ├── repository            # Database access layer
│   └── service               # Business logic layer
├── pkg
│   └── utils                 # Shared utilities (Password hashing, JWT)
├── docs                      # Generated Swagger docs
└── .env                      # Environment variables
```

## ⚙️ Setup & Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.20+
- [MongoDB](https://www.mongodb.com/try/download/community) installed and running locally or a cloud instance (Atlas).

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/magicstreammovies-server.git
cd magicstreammovies-server
```

### 2. Environment Configuration

Create a `.env` file in the root directory:

```env
MONGODB_URL=mongodb://localhost:27017
DATABASE_NAME=MagicStreamMovies
SECRET_KEY=your_secret_key_here
SECRET_REFRESH_KEY=your_refresh_secret_key_here
OPEN_API_KEY=your_openai_api_key
BASE_PROMPT_TEMPLATE="Rate the sentiment of this review based on the following rankings: {rankings}"
RECOMMENDED_MOVIE_LIMIT=5
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Install Dependencies

```bash
go mod download
```

## 🏃‍♂️ Running the Application

### Standard Run

```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:8080`.

### Live Reload (Recommended for Dev)

If you have [Air](https://github.com/air-verse/air) installed:

```bash
air
```

## 🧪 Testing

Run unit tests using the standard Go test command:

```bash
go test ./internal/...
```

## 📖 API Documentation

The API is documented using Swagger. Once the server is running, visit:

👉 **http://localhost:8080/swagger/index.html**

### Authentication in Swagger

1. Login via `/login` to get a token.
2. Click the **Authorize** button in Swagger.
3. Enter your token (e.g., `eyJhbGci...`). You do **not** need to type "Bearer ".

## 🤝 Contributing

1. Fork the project.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
