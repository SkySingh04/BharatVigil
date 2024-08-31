# Centralized Application-Context Aware Firewall

## Overview

This project is a centralized application-context aware firewall designed to manage and monitor network traffic for applications running on endpoint devices. The firewall allows granular control over domains, IP addresses, and protocols for each application and can be managed through a centralized web console. It also includes AI/ML capabilities to detect and alert on abnormal network behavior.

## Features

- **Granular Application Control:** Restrict or allow specific domains, IPs, and protocols for each application.
- **Centralized Management:** Manage firewall rules and monitor network activity from a centralized web console.
- **AI/ML Integration:** Detect and alert on anomalous network behavior using AI/ML models.
- **Cross-Platform Support:** Designed to work on both Windows and Linux endpoints.

## Directory Structure

```plaintext
project-root/firewall-tool
│
├── cmd/                    # Cobra commands for the CLI
│   ├── reload.go           # Command to reload the config file
│   └── root.go             # Root command and entry point for the CLI
│
├── config/                 # Configuration handling
│   └── config.go           # Go structs and methods to load and parse the config file
│
├── pkg/                    # Core logic for the firewall tool
│   ├── firewall/           # Application-specific firewall logic
│   ├── monitoring/         # Network monitoring logic
│   ├── ai_ml/              # AI/ML integration logic
│   └── web_console/        # Web console server logic
│
├── config.yaml             # Configuration file for the firewall
├── main.go                 # Main entry point for the CLI tool
├── README.md               # Project documentation
└── go.mod                  # Go module file
```

## Configuration File (`config.yaml`)

The `config.yaml` file is the core configuration file for the firewall tool. It includes various settings that control the behavior of the firewall, monitoring, AI/ML detection, and the web console. Below is an overview of the key sections:

### **Key Sections:**

- **firewall:** Defines rules for each application, specifying allowed and blocked domains, IP addresses, and protocols.
- **monitoring:** Configures network traffic monitoring, including logging and alert thresholds.
- **ai_ml:** Configures AI/ML model integration, including the model endpoint and anomaly detection settings.
- **endpoints:** Lists all endpoints where the firewall agent is deployed, including their OS, IP address, and hostname.
- **web_console:** Configures the web console, including port settings, access control via allowed/blocked IPs, and admin users.
- **logging:** Manages log settings, including log level, file paths, and rotation policies.
- **network:** Configures deep packet inspection, TLS certificate paths, and proxy settings.

## ML Team Tasks

The AI/ML team is responsible for developing and deploying the model that will detect anomalous network behavior. Here are the key tasks:

1. **Model Development:**
   - Develop an AI/ML model capable of detecting abnormal network traffic patterns.
   - Train the model using relevant network traffic data.
   
2. **Model Deployment:**
   - Deploy the model as a service, either locally or on a remote server.
   - Ensure the model is accessible via an API endpoint (e.g., `http://localhost:5000/predict`).

3. **Integration:**
   - Work with the backend team to ensure the firewall tool can communicate with the model endpoint.
   - Define the input/output schema for the model to facilitate smooth integration.

4. **Testing:**
   - Test the model's accuracy and responsiveness when integrated with the firewall tool.
   - Provide feedback to the backend team for any necessary adjustments.

## Frontend Team Tasks

The frontend team is responsible for designing and developing the centralized web console. Key tasks include:

1. **UI/UX Design:**
   - Create a user-friendly interface that allows admins to manage firewall rules, view logs, and monitor network activity.
   - Use design inspirations such as [Portmaster](https://safing.io/) for UI elements.

2. **Web Console Development:**
   - Develop the web console using a modern frontend framework (e.g., React, Svelte).
   - Implement features to:
     - View and edit the firewall configuration.
     - Monitor real-time network traffic and logs.
     - Display alerts for abnormal network behavior.

3. **API Integration:**
   - Integrate with the backend API to load and save configurations.
   - Provide real-time updates and control over the firewall tool through the web console.

4. **Testing:**
   - Ensure the web console is fully functional and integrates seamlessly with the backend tool.
   - Perform usability testing and gather feedback for improvements.

## Getting Started

To get started with the project:

1. **Clone the Repository:**
   ```bash
   git clone https://github.com/your-username/centralized-firewall.git
   cd centralized-firewall
   ```

2. **Install Dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run the CLI Tool:**
   ```bash
   go run main.go
   ```

4. **Reload Configurations:**
   ```bash
   go run main.go reload-config --config config.yaml
   ```
