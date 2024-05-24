# Go TCP Server-Client Handshake

This project is a simple implementation of a TCP server-client handshake in Go.

## Project Structure

The main file of the project is `cmd/main.go`. This file contains the main function which starts the server and client, 
and the functions for the server and client operations.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing 
purposes.

### Prerequisites

- Go 1.16 or higher

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/tanryberdi/network-handshake.git
```

### Running the project
I used niktrix/btcd-docker for bitcoin node implementation. Be sure for running this image and using port 8333 for 
bitcoin node.

```bash
go run cmd/main.go
```

