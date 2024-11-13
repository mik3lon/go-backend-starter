# (WIP) Go HTMX Starter

![Go](https://img.shields.io/badge/Go-1.18+-blue.svg)
![Gin](https://img.shields.io/badge/Gin-Framework-blue.svg)
![HTMX](https://img.shields.io/badge/HTMX-JS-orange.svg)
![TailwindCSS](https://img.shields.io/badge/TailwindCSS-v3.0-green.svg)
![Hexagonal Architecture](https://img.shields.io/badge/Architecture-Hexagonal-blueviolet.svg)

Welcome to **Go HTMX Starter** — a robust starter kit designed to build scalable and interactive web applications using **Go** with **Gin**, **HTMX**, and **TailwindCSS** in a **Hexagonal Architecture** pattern. This project aims to provide a clean, organized codebase to kickstart your Go projects with modern frontend capabilities.

## 🚀 Features

- **Gin**: A lightweight, high-performance HTTP web framework for Go.
- **Hexagonal Architecture**: Also known as "Ports and Adapters," this architecture provides a structured approach for better separation of concerns and testability.
- **HTMX**: Enhances the interactivity of the application by handling AJAX requests and updating parts of the page without reloading.
- **TailwindCSS**: A utility-first CSS framework for fast and responsive UI development.
- **Ready-to-Use Structure**: Boilerplate code to help you focus on writing your application logic.

## 🛠️ Technologies Used

| Technology           | Purpose                                        |
|----------------------|------------------------------------------------|
| **Go**               | The main programming language.                 |
| **Gin**              | Web framework for routing and HTTP handling.   |
| **Hexagonal Architecture** | For organizing code and business logic.   |
| **HTMX**             | Enables interactive web features without JavaScript complexities. |
| **TailwindCSS**      | For styling and creating a modern UI.          |

---

## ⚙️ Getting Started

### Prerequisites

Ensure you have the following installed:

- [Go 1.18+](https://golang.org/doc/go1.18)
- [Node.js and npm](https://nodejs.org/) (for TailwindCSS compilation)
- [Go Gin](https://github.com/gin-gonic/gin)

### Installation

1. **Clone the repository**

    ```bash
    git clone https://github.com/mik3lon/starter-template.git
    cd starter-template
    ```

2. **Install Go dependencies**

    ```bash
    go mod tidy
    ```

3. **Install TailwindCSS**

   You can install Tailwind via npm or yarn:

    ```bash
    npm install
    ```

4. **Compile TailwindCSS**

   Run the following command to watch and compile your TailwindCSS assets:

    ```bash
    npm run build:css
    ```

5. **Run the server**
   Set the .env with the params and start the Go Gin server with templ:

    ```bash
    templ generate --watch --proxy="http://localhost:8000" --cmd="go run cmd/web/app/main.go"
    ```

6. **Access the application**

   Open your browser and go to [http://localhost:8000](http://localhost:8000) to view the app.

---

## ✨ Key Components

### Gin

Gin is a high-performance HTTP web framework for Go that simplifies building web applications. It's used here to handle routing, middleware, and HTTP request/response processing.

### Hexagonal Architecture

This architecture helps separate the core business logic (domain) from infrastructure and application services. It promotes maintainability and testability by keeping the application logic free from dependencies on frameworks and external tools.

### HTMX

HTMX enables AJAX, CSS Transitions, WebSockets, and Server-Sent Events (SSE) directly in HTML, making it easier to build interactive components without writing custom JavaScript. It’s ideal for scenarios where you want dynamic, responsive behavior without a heavy frontend framework.

### TailwindCSS

TailwindCSS is a utility-first CSS framework that allows for rapid UI design. This project uses Tailwind to style components and create a modern, responsive UI.


# Makefile Commands

The `Makefile` provides a convenient way to manage your Docker containers and run tests. Below are the available targets:

## Usage

```bash
make [target]
```

| Target | Description                                                        |
|--------|--------------------------------------------------------------------|
| `up`   | Start the application containers in detached mode.                 |
| `down` | Stop the application containers and remove the network.            |
| `logs` | View logs of the running containers in real-time.                  |
| `clean`| Stop the containers, remove volumes, and orphan containers.        |
| `tests`| Execute tests in the project and display the result.               |
