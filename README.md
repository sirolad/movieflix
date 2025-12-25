# MagicStreamMovies

MagicStreamMovies is a full-stack movie streaming application built with Go (Golang) and React. It allows users to browse movies, get recommendations, and manage their favorite genres.

## Tech Stack

### Backend
*   **Language:** Go (Golang)
*   **Framework:** Gin Web Framework
*   **Database:** MongoDB
*   **Authentication:** JWT (JSON Web Tokens)
*   **Documentation:** Swagger
*   **Other Libraries:** `godotenv`, `validator`, `crypto`

### Frontend
*   **Framework:** React
*   **Build Tool:** Vite
*   **Styling:** Bootstrap (React Bootstrap)
*   **Routing:** React Router
*   **HTTP Client:** Axios

## Project Structure

```
MagicStreamMovies/
├── Client/
│   └── magic-stream-client/  # React Frontend
└── Server/
    └── MagicStreamMoviesServer/ # Go Backend
```

## Getting Started

### Prerequisites
*   Go 1.25+
*   Node.js & npm
*   MongoDB instance

### Backend Setup
1.  Navigate to the server directory:
    ```bash
    cd Server/MagicStreamMoviesServer
    ```
2.  Install dependencies:
    ```bash
    go mod download
    ```
3.  Create a `.env` file with your configuration (MongoDB URI, JWT Secret, etc.).
4.  Run the server:
    ```bash
    go run main.go
    ```

### Frontend Setup
1.  Navigate to the client directory:
    ```bash
    cd Client/magic-stream-client
    ```
2.  Install dependencies:
    ```bash
    npm install
    ```
3.  Start the development server:
    ```bash
    npm run dev
    ```

## Features
*   User Authentication (Login/Register)
*   Browse Movies
*   Personalized Movie Recommendations
*   Responsive UI
