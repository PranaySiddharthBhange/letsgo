# Lets' Go - Command Line Messaging Platform

Lets' Go is a simple command line messaging platform written in Go. It allows you to send and receive messages through the terminal.

## Installation

### Option 1: Manual Installation

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/PranaySiddharthBhange/letsgo.git
   ```

2. Navigate to the project directory:

   ```bash
   cd letsgo
   ```

3. Run the main application:
   ```bash
   go run main.go
   ```

### Option 2: Install as a Custom Command

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/PranaySiddharthBhange/letsgo.git
   ```

2. Navigate to the project directory:

   ```bash
   cd letsgo
   ```

3. Give execute permission to the `letsgo.sh` script:

   ```bash
   chmod +x letsgo.sh
   ```

4. Run the script to set up the custom command:

   ```bash
   ./letsgo.sh
   ```

   After this, you can use the custom command by typing `letsgo` in the terminal.

### Option 3: Docker Installation

1. Pull the image from Docker Hub:

   ```bash
   docker pull pranaybhange/letsgo
   ```

2. Run the image:

   ```bash
   docker run -i pranaybhange/letsgo
   ```

## Features

1. **Create an Account**

   - Use Lets' Go to create a user account with a unique username and password.

2. **Send Messages**

   - Send messages to other users through the terminal.

3. **Receive Messages**
   - Receive and view messages sent to your account.
4. **Ping Me**
   - username : pranay

## Prerequisites

Ensure that you have Git and Go installed on your device before using Option 1 and 2.

```bash
# Install Git
sudo apt-get update
sudo apt-get install git

# Install Go
# Follow the official Go installation guide: https://golang.org/doc/install
```

**Note:** Message and password encryption are not implemented yet. Exercise caution when using this tool for sensitive information.

Feel free to contribute to the project by submitting issues or pull requests.
