# URL Shortener DevSecOps Project

## 1. Application Overview

The URL Shortener Service is a high-performance web application that converts long URLs into compact, shareable codes. It emphasizes security and scalability through a DevSecOps approach.

### Architecture Diagram
The following diagram illustrates the complete DevSecOps flow, from code commit to production deployment.

![Architecture Diagram](architecture.png)

## 2. How to Run Locally

### Prerequisites
- **Go**: 1.24 or higher
- **Docker**: For containerization

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/souvik-cyclic/url-shortener-devops.git
   cd url-shortener-devops
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Run the application**
   ```bash
   go run main.go
   ```
   The server will start at `http://localhost:8080`.

### Running with Docker

1. **Build the image**
   ```bash
   docker build -t url-shortener .
   ```

2. **Run the container**
   ```bash
   docker run -p 8080:8080 url-shortener
   ```

## 3. Secrets Configuration

To enable the automated CI/CD pipeline, you must configure the following **GitHub Actions Secrets** in your repository settings:

| Secret Name | Description | Required For |
| :--- | :--- | :--- |
| `DOCKERHUB_USERNAME` | Your Docker Hub username. | Pushing Docker images in CI. |
| `DOCKERHUB_TOKEN` | Docker Hub Access Token. | Authenticating with Docker Hub. |

> **Note:** The application itself (`main.go`) runs with default settings locally and does not require a `.env` file for basic operation.

## 4. CI/CD Pipeline Explanation

The project utilizes GitHub Actions for both Continuous Integration (CI) and Continuous Deployment (CD).

### **CI Pipeline** (Triggers on `push` to `main`)
1. **Linting**: Uses `golangci-lint` to enforce code quality and style.
2. **Security - SAST**: `CodeQL` analyzes source code for vulnerabilities.
3. **Security - SCA**: `Go Vulncheck` scans dependencies for known issues.
4. **Unit Testing**: Runs Go tests to ensure logic correctness.
5. **Build**: Compiles the Go binary.
6. **Docker Build**: Creates the container image.
7. **Security - Container Scan**: `Trivy` scans the Docker image for Critical/High vulnerabilities.
8. **Push**: Uploads the secure image to Docker Hub.

### **CD Pipeline** (Triggers on successful CI)

> **Requirement:** Deployment requires a **Self-Hosted Runner** configured for this repository, which must have access to a **Kubernetes cluster**.

1. **Deploy**: Connects to the self-hosted Kubernetes runner.
2. **Update**: Applies Kubernetes manifests (`deployment.yaml`, `service.yaml`).
3. **Verify**: Checks rollout status and performs connectivity tests.
