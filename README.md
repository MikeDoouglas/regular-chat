# [RegularChat.com](https://regularchat.unibutton.com/)
### Just a regular chat!

![man typing](https://uwaterloo.ca/writing-and-communication-centre/sites/default/files/uploads/images/email_wink.gif)

- [Overview](#overview)
- [Technologies Used](#technologies-used)
- [Repository Structure](#repository-structure)
- [Usage](#usage)

## Overview

This project is a simple chat application with a ReactJS frontend and a Golang backend communicating through WebSocket.

## Technologies Used

- **Frontend**: ReactJS
  - `create-react-app`
- **Backend**: Go (Golang) 1.22

## Usage

**Start the backend**:
1. Duplicate `.env.example` and rename the copy to `.env`
2. Fill in the environment variables
3. Update dependencies

    ```sh
    cd backend/
    go mod tidy
    ```
4. Now you can run using:

    ```sh
    go run cmd/chat/main.go
    ```
    or
    ```sh
    docker build -t chat-backend .
    docker run -d --name chat-backend -p 8080:8080 -e ENV=develop -e PORT=8080 chat-backend
    ```

**Start the frontend**:
1. Just install dependencies and run:
    ```sh
    cd frontend
    npm install
    npm start
    ```
    or 

    ```sh
    docker build -t chat-frontend .
    docker run -d --name chat-frontend -p 3000:3000 chat-frontend
    ```

**Acess chat**

Open in your browser `http://localhost:3000`.

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
