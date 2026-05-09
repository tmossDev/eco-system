# Metric Boat Database

## UML Diagram
*To do: Add UML Diagram*

## Local Setup
To set up the database locally, follow these steps:

1. **Install Docker:**
   - **Windows/Mac:**
     1. Download and install Docker Desktop from the official Docker website.
     2. Follow the installation instructions specific to your operating system.
     3. Once installed, start Docker Desktop.

   - **Linux:**
     1. Update your package database:
        ```bash
        sudo apt-get update
        ```
     2. Install Docker using the convenience script provided by Docker:
        ```bash
        curl -fsSL https://get.docker.com -o get-docker.sh
        sudo sh get-docker.sh
        ```
     3. Add your user to the Docker group to run Docker commands without `sudo`:
        ```bash
        sudo usermod -aG docker $USER
        ```
     4. Log out and log back in to apply the group changes.

2. **Start DB container:**
   ```bash
   docker compose up --build
   ```
