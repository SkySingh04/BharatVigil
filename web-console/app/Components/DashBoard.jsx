"use client";
import React, { useEffect, useState } from "react";
import PostCard from "./PostCard";
import MonitoringCard from "./MonitoringCard";
import AiMlCard from "./AiMlCard";
import LoggingCard from "./LoggingCard";

import yaml from "js-yaml";
import { FaBell, FaExclamationTriangle, FaCog, FaSignOutAlt } from "react-icons/fa";

const DashBoard = () => {
  const [rules, setRules] = useState([]);
  const [monitoring, setMonitoring] = useState({}); // State for monitoring configuration
  const [AiMl, setAiMl] = useState({}); // State for AI/ML configuration
  const[logging, setLogging] = useState({}); // State for logging configuration
  const[endpoints, setEndpoints] = useState([]); // State for endpoints
  const [requests, setRequests] = useState([]); // State for requests
  const [loadingRequests, setLoadingRequests] = useState(true); // Loading state for requests
  const [isRequestScreenOpen, setIsRequestScreenOpen] = useState(false); // State for request screen
  

  // Fetch firewall rules from backend
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch("http://localhost:8080/config");
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);

        // Access firewall rules from fetched data
        const firewallRules = data.firewall?.rules || [];
        setRules(firewallRules);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };

    fetchConfig();
  }, []);

  //Fetch monitoring data from backend
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch("http://localhost:8080/config");
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);
  
        // Access monitoring configuration from fetched data
        const monitoringConfig = data.monitoring || {};
        setMonitoring(monitoringConfig);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };
  
    fetchConfig();
  }, []);

  // Fetch AI/ML data from backend
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch("http://localhost:8080/config");
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);
  
        // Access monitoring configuration from fetched data
        const AiMlConfig = data.ai_ml || {};
        setAiMl(AiMlConfig);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };
  
    fetchConfig();
  }, []);

  // Fetch logging data from backend
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch("http://localhost:8080/config");
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);
  
        // Access logging configuration from fetched data
        const loggingConfig = data.logging || {};
        console.log("Fetched logging config:", loggingConfig); // Debug log
        setLogging(loggingConfig);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };
  
    fetchConfig();
  }, []);

  //Fetch Endpoint data from backend
  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const response = await fetch("http://localhost:8080/config");
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);
  
        // Access logging configuration from fetched data
        const endpointsConfig = data.endpoints || {};
        console.log("Fetched logging config:", endpointsConfig); // Debug log
        setEndpoints(endpointsConfig);
      } catch (error) {
        console.error("Error fetching config:", error);
      }
    };
  
    fetchConfig();
  }, []);

  // Reusable function to fetch the list of requests (from /requests)
  const fetchRequests = async () => {
    try {
      const response = await fetch("http://localhost:8080/requests");
      if (!response.ok) {
        throw new Error("Failed to fetch requests");
      }
      const data = await response.json();
      setRequests(data);
      setLoadingRequests(false); 
    } catch (error) {
      console.error("Error fetching requests:", error);
      setLoadingRequests(false); // Stop loading on error
    }
  };

  // Fetch initial requests on component load
  useEffect(() => {
    fetchRequests();
  }, []);

  // Use SSE to fetch real-time events from the backend (using /events endpoint)
  useEffect(() => {
    const eventSource = new EventSource("http://localhost:8080/events");
    console.log("SSE connection opened", eventSource);

    eventSource.onmessage = (event) => {
      console.log("New SSE event:", event.data); // Log event data
      const data = JSON.parse(event.data);
      setRequests((prevRequests) => [...prevRequests, data]);
    };

    eventSource.onerror = (error) => {
      console.error("SSE Error:", error); // Log error
      eventSource.close();
    };

    return () => {
      eventSource.close(); // Clean up
    };
  }, []);

  const handleRequestOpen = () => {
    setIsRequestScreenOpen(true);
  };

  // Function to handle closing the screen
  const handleRequestClose = () => {
    setIsRequestScreenOpen(false);
  };

  return (
    <div className="bg-black grid grid-cols-4 gap-4 mt-6">
    {/* Left Section */}
    <div className="w-8 bg-black p-2 flex flex-col items-center h-full fixed left-0 top-0 ">
    {/* Icon Buttons */}
    <div className="flex-1 flex flex-col items-center gap-4">
      <FaBell size={24} className="hover:text-green-400 cursor-pointer" />
      <FaExclamationTriangle size={24} className="hover:text-green-400 cursor-pointer" />
      <FaCog size={24} className="hover:text-green-400 cursor-pointer" />
    </div>
    {/* Logout Button */}
    <div className="mt-auto">
      <FaSignOutAlt size={24} className="hover:text-red-400 cursor-pointer" />
    </div>
  </div>

  {/* Main Content Section */}
  <div className=" w-full p-4 flex flex-col bg-black text-white">
    <div className="flex-1 overflow-y-auto">
      {/* Profile Section */}
      <div className="flex flex-col items-center justify-center mb-4">
        <img src="/image.png" className="w-40 h-40 mb-2 rounded-full" alt="Secure" />
        <strong className="text-white mt-2">Connection Secure</strong>
      </div>

      {/* Search Input */}
      <div className="mb-4">
        <label className="input input-bordered flex items-center gap-2">
          <input type="text" className="grow bg-black text-white border-gray-500 rounded-lg p-2" placeholder="Search" />
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" className="h-4 w-4 opacity-70">
            <path fillRule="evenodd" d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z" clipRule="evenodd" />
          </svg>
        </label>
      </div>

      {/* Firewall Rules Section */}
      <div className="bg-base-200 rounded-lg p-4 mb-4 border border-white border-opacity-20 h-[250px] overflow-y-auto">
        <ul className="space-y-4">
          {['Microsoft Edge', 'System DNS Client', 'Firefox', 'Chrome', 'Safari'].map((app, index) => (
            <li key={index} className="p-4 bg-gray-800 rounded-lg border border-gray-600 flex justify-between items-center">
              <strong className="text-white">{app}</strong>
              <div className="radial-progress bg-primary text-primary-content border-primary" style={{ "--value": "50", "--size": "3rem" }} role="progressbar">
                50%
              </div>
            </li>
          ))}
        </ul>
      </div>

      {/* Connections Information */}
      <div className="bg-base-200 rounded-lg p-4 border border-white border-opacity-20 h-[250px] overflow-y-auto">
        <ul className="space-y-4">
          {[
            { country: "Canada", ip: "78.138.17.182", connections: 2, hops: 2 },
            { country: "USA", ip: "5.34.178.198", connections: 45, hops: 2 },
            { country: "France", ip: "141.95.158.73", connections: 6, hops: 2 },
            { country: "Germany", ip: "138.201.140.70", connections: 184, hops: 2 },
          ].map((data, index) => (
            <li key={index} className="p-4 bg-gray-800 rounded-lg border border-gray-600 flex items-center">
              <img src="/path-to-your-icon.png" alt="Country Icon" className="w-6 h-6 mr-4" />
              <div>
                <h3 className="text-white text-lg font-semibold">{data.ip}</h3>
                <ul className="text-sm text-gray-400 space-y-1">
                  <li><strong>Country:</strong> {data.country}</li>
                  <li><strong>Connections:</strong> {data.connections}</li>
                  <li><strong>HOPS:</strong> {data.hops}</li>
                </ul>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  </div>

  {/* Middle Section */}
  <div className="col-span-3 p-6 flex flex-col bg-black h-screen">
  <div className="h-full rounded-lg flex flex-col justify-between relative">
    <div className="grid grid-cols-2 gap-4 h-full">
      {/* Left Section: PostCard with Scroll */}
      <div className="h-full overflow-y-auto rounded-lg">
        <div className="space-y-4">
          {rules.map((rule) => (
            <PostCard
              key={rule.id}
              application={rule.application}
              allowedDomains={rule.allowed_domains}
              blockedDomains={rule.blocked_domains}
              allowedIps={rule.allowed_ips}
              blockedIps={rule.blocked_ips}
              protocols={rule.protocols}
            />
          ))}
        </div>
      </div>

      {/* Right Section: MonitoringCard, AiMlCard, LoggingCard */}
      <div className="flex flex-col space-y-4 overflow-y-auto rounded-lg relative">
        <MonitoringCard
          enable={monitoring.enable}
          logFile={monitoring.log_file}
          alertThresholds={monitoring.alert_thresholds}
        />
        <AiMlCard
          modelEndpoint={AiMl.model_endpoint}
          enableAnomaly={AiMl.enable_anomaly_detection}
        />
        <LoggingCard
          logLevel={logging.log_level}
          logFile={logging.log_file}
          maxSize={logging.max_size}
          maxBackups={logging.max_backups}
          maxAge={logging.max_age}
        />
        {/* Requests Button */}
        <div className="flex justify-center">
          <button 
            className="btn btn-info btn-outline mt-4"
            onClick={handleRequestOpen}
          >
            Requests
          </button>
        </div>

        {isRequestScreenOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-60 flex items-center justify-center">
          <div className="bg-gray-800 p-6 rounded-lg max-w-lg w-full">
            {/* Loading State */}
            {loadingRequests ? (
              <div className="text-white loading loading-infinity loading-lg">Loading requests</div>
            ) : requests.length > 0 ? (
              <div className="overflow-y-auto w-full h-80">
                {/* Request List */}
                {requests.slice().reverse().map((request) => (
                  <div key={request.no} className="bg-base-100 hover:bg-base-300 border-b border-white border-opacity-20 text-white p-4">
                    <div className="flex flex-wrap items-center gap-4">
                      <span className="text-sm">
                        <strong>Time:</strong> {request.time}
                      </span>
                      <span className="text-sm">
                        <strong>Source:</strong> {request.source} {'-->'} <strong>Destination:</strong> {request.destination}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-white">No requests found.</div>
            )}

            {/* Close Button */}
            <div className="flex justify-end mt-4">
              <button 
                className="btn btn-error btn-outline"
                onClick={handleRequestClose}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
      </div>
    </div>
  </div>
</div>

</div>
  );
};

export default DashBoard;
