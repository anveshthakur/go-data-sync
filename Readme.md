# Go-Sync-Data

`Go-Sync-Data` is a fullstack application that can be built, run, and managed using the included Makefile and Dockerfile. This project includes both backend and UI components.

## 🚀 Features

- Real-time database synchronization
- Interactive table view and data comparison
- Docker-based PostgreSQL setup

## Prerequisites
- **Go** 
- **Node.js & npm**
- **Docker**
- **Make** (for running Makefile commands)

## Screenshots

### HomePage
![Home](./assets/home.png)

### Table View Page
![Table View Page](./assets/table-view.png)


---


## Installation
### 1. Clone the repository

```sh
git clone https://github.com/anveshthakur/go-data-sync
```

### 2. Database Setup 
```sh

// create an Image from the dockerfile
docker build -t go-sync-data .

// run the created docker image 
docker run -d --name postgres-container -p 5432:5432 go-sync-data

```

### 3.Environment Configuration
```
// UI env file
NEXT_PUBLIC_BACKEND_API=http://localhost:8080
```

## 🎯 Usage

### Using Makefile Commands

The Makefile provides commands to manage the application. Below is a list of commands with descriptions.

### Build the Application

To compile the `syncData` backend application:

```bash
make build
```


### Run the Application


```bash
make start
```

### Stop the Application


```bash
make stop
```

### Install Dependencies in UI

```bash
make install_ui
```

### Run the UI

To compile the `syncData` backend application:

```bash
make start_ui
```

## 📊 Database Connection Details

### Source Database
```
host: localhost
port: 5432
username: postgres
password: password
dbname: sourcedb
```

### Target Database
```
host: localhost
port: 5432
username: postgres
password: password
dbname: targetdb
```

## File Structure
```
./go-data-sync
├── assets
├── bin
├── cmd
│   ├── tmp
│   └── web
├── scripts
├── tmp
└── ui
    ├── app
    │   └── dashboard
    ├── components
    │   └── ui
    ├── hooks
    ├── lib
    └── tmp
```