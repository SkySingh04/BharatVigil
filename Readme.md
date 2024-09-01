# BharatVigil - Centralized Application-Context Aware Firewall

## Overview

BharatVigil is a centralized application-context aware firewall designed to manage and monitor network traffic for applications running on endpoint devices. The firewall allows granular control over domains, IP addresses, and protocols for each application and can be managed through a centralized web console. It also includes AI/ML capabilities to detect and alert on abnormal network behavior.

## TODO:
- Implement blocking from config file
- Implement parsing of all network data
- Implement config file filtering from web console.
- Implement one-click install and packaging (needs to be part of ppt like keploy)
- Implement ML model to provide context awareness
- Enable monitoring for all systems connected on the network in a centralized way.
- Add Priority level for applications and add priority level for network computers.
- Future: Add protection for sharing within networks.

## Features

- **Granular Application Control:** Restrict or allow specific domains, IPs, and protocols for each application.
- **Centralized Management:** Manage firewall rules and monitor network activity from a centralized web console.
- **AI/ML Integration:** Detect and alert on anomalous network behavior using AI/ML models.
- **Cross-Platform Support:** Designed to work on both Windows and Linux endpoints.

### API Documentation

#### Base URL
```
http://localhost:8080
```

#### Endpoints

---

**1. Get Current Configuration**

- **URL**: `/config`
- **Method**: `GET`
- **Description**: Retrieves the current YAML configuration file.
- **Response**:
  - **200 OK**: Returns the YAML file content.
  - **404 Not Found**: If the configuration file is not found.

---

**2. Update Configuration**

- **URL**: `/config`
- **Method**: `POST`
- **Description**: Updates the YAML configuration file with new data.
- **Request Body**:
  - **Content-Type**: `application/json`
  - **Body**: 
    ```json
    {
      "key1": "value1",
      "key2": "value2"
    }
    ```
- **Response**:
  - **200 OK**: Configuration updated successfully.
  - **400 Bad Request**: Invalid request format.
  - **500 Internal Server Error**: Failed to write the configuration file.

---

**3. Network Activity (Placeholder)**

- **URL**: `/network-activity`
- **Method**: `GET`
- **Description**: Placeholder endpoint for sending network activity data.
- **Response**:
  - **200 OK**: Indicates network activity sending.


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


## Flags

BharatVigil includes several flags that allow you to control its behavior when starting the tool:

- **`--disable-console`**: Disable the web console. Use this flag if you want to start the tool without launching the web console.
  ```bash
  go run main.go --disable-console
  ```

- **`--disable-server`**: Disable the backend server. Use this flag if you want to start the tool without launching the backend server.
  ```bash
  go run main.go --disable-server
  ```

- **`--disable-ml`**: Disable the ML model. Use this flag if you want to start the tool without launching the ML model.
  ```bash
  go run main.go --disable-ml
  ```

- **`--disable-ebpf`**: Disable the eBPF program. Use this flag if you want to start the tool without loading the eBPF program.
  ```bash
  go run main.go --disable-ebpf
  ```

- **`--config`**: Specify the path to the configuration file. If not provided, the tool will use the default `config.yaml` in the root directory.
  ```bash
  go run main.go --config /path/to/config.yaml
  ```

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
