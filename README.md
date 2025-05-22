# [RegularChat.com](https://regularchat.unibutton.com/)
### Just a regular chat!

![man typing](https://uwaterloo.ca/writing-and-communication-centre/sites/default/files/uploads/images/email_wink.gif)

- [Overview](#overview)
- [Technologies Used](#technologies-used)
- [Repository Structure](#repository-structure)
- [Usage](#usage)

## Overview

RegularChat is a lightweight chat application featuring a ReactJS frontend and a Golang backend, communicating over WebSocket.

## Technologies Used

- **Frontend**: ReactJS  
- **Backend**: Go (Golang) 1.22

## Usage

### Starting the Backend

1. Copy the `.env.example` file and rename it to `.env`.
2. Fill in the required environment variables.
3. Install dependencies:

    ```sh
    cd backend/
    go mod tidy
    ```

4. Start the backend:

    ```sh
    go run cmd/chat/main.go
    ```

    Or use Docker:

    ```sh
    docker build -t chat-backend .
    docker run -d --name chat-backend -p 8080:8080 -e ENV=develop -e PORT=8080 chat-backend
    ```

### Starting the Frontend

1. Install dependencies and start the app:

    ```sh
    cd frontend
    npm install
    npm start
    ```

    Or use Docker:

    ```sh
    docker build -t chat-frontend .
    docker run -d --name chat-frontend -p 3000:3000 chat-frontend
    ```

### Accessing the Chat

Open your browser and navigate to:  
**`http://localhost:3000`**

## Repository Structure

```
├── backend
│   ├── cmd
│   │   └── chat
│   │       └── main.go
│   ├── Dockerfile
│   ├── files
│   │   ├── adjectives.txt
│   │   ├── celebrities.txt
│   │   ├── fantasies.txt
│   │   ├── foods.txt
│   │   ├── objects.txt
│   │   └── profissions.txt
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       └── chat
│           ├── handler.go
│           ├── models.go
│           ├── name_generator.go
│           ├── room.go
│           └── user.go
├── frontend
│   ├── Dockerfile
│   ├── package.json
│   ├── package-lock.json
│   ├── public
│   │   ├── favicon.ico
│   │   ├── index.html
│   │   ├── logo.png
│   │   ├── manifest.json
│   │   └── robots.txt
│   └── src
│       ├── App.css
│       ├── App.js
│       ├── Header.css
│       ├── Header.js
│       ├── index.css
│       ├── index.js
│       └── reportWebVitals.js
└── README.md
```
