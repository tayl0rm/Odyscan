# Odyscan

Odyscan is a container image security scanning tool designed to pull images from Google Cloud's Artifact Registry, extract their filesystem, and scan for malware using ClamAV. The tool is built to run on a Kubernetes cluster (GKE, K3s, etc.) and provides a web-based interface for easy interaction.

## Features

- **Pull images from Google Artifact Registry**
- **Extract image files for scanning**
- **Scan extracted files using ClamAV**

## Prerequisites

- Kubernetes cluster (K3s recommended)
- Google Artifact Registry Repository
- ClamAV deployment in Kubernetes
- Go 1.23.4
- Node.js (for frontend modifications)

## Installation

### Clone Repository

```sh
git clone https://github.com/tayl0rm/odyscan.git
cd odyscan
```

## Running Locally

To run the application locally for testing:

```sh
go run main.go
```

Access the UI at `http://localhost:8080`.

## Usage

1. Enter the container image name in the web UI.
2. Click "Scan" to start the process.
3. View real-time scan progress.
4. Results will indicate whether the image contains threats.

## Running Tests

Unit tests can be executed with:

```sh
cd odyscan/scanner
go test 
```

## License

[MIT License](LICENSE)
