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

  // Handler to refresh the requests list manually
  const handleRefresh = () => {
    setLoadingRequests(true); // Show loading while fetching
    fetchRequests(); // Re-fetch requests from the backend
  };

  return (
    <div className="min-h-screen bg-black grid grid-cols-4 gap-4">
  {/* Left Section */}
  <div className="col-span-1">
    <div className="flex h-screen bg-black text-white">
      <div className="w-16 bg-black-900 p-4 flex flex-col items-center">
        {/* Icon buttons */}
        <div className="mb-8">
          <FaBell size={24} className="mb-8 hover:text-green-400 cursor-pointer" />
          <FaExclamationTriangle size={24} className="mb-8 hover:text-green-400 cursor-pointer" />
          <FaCog size={24} className="hover:text-green-400 cursor-pointer" />
        </div>
        <div className="mt-auto">
          <FaSignOutAlt size={24} className="hover:text-red-400 cursor-pointer" />
        </div>
      </div>

      <div className="bg-black p-2">
        <div className="flex flex-col items-center justify-center mb-4">
          <img src="/image.png" className="w-40 h-40 mb-2 rounded-full" alt="Secure" />
          <strong className="text-white mt-2">Connection Secure</strong>
        </div>

        <div className="mb-4">
          <label className="input input-bordered flex items-center gap-2">
            <input type="text" className="grow" placeholder="Search" />
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" fill="currentColor" className="h-4 w-4 opacity-70">
              <path fillRule="evenodd" d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z" clipRule="evenodd" />
            </svg>
          </label>
        </div>

        <div className="menu bg-base-200 rounded-box w-100 mb-4 overflow-y-auto h-80 border border-white border-opacity-20">
          <ul>
            {['Microsoft Edge', 'System DNS Client', 'Firefox'].map((app, index) => (
              <li key={index} className="p-4 bg-base-100 rounded-lg mb-2 border border-white border-opacity-20">
                <div className="flex justify-between items-center">
                  
                  <strong className="text font-semibold">{app}</strong>
                  <div className="radial-progress bg-primary text-primary-content border-primary" style={{  "--value": "50", "--size": "5rem"}}role="progressbar">
                    50%
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>

        {/* Connections Information */}
        <div className="menu bg-base-200 rounded-box w-100 mt-4 h-80 border border-white border-opacity-20">
          <div className="h-64 overflow-y-auto">
            <ul className="grid grid-cols-1 gap-4">
              {[
                { country: "Canada", ip: "78.138.17.182", connections: 2, hops: 2 },
                { country: "USA", ip: "5.34.178.198", connections: 45, hops: 2 },
                { country: "France", ip: "141.95.158.73", connections: 6, hops: 2 },
                { country: "Germany", ip: "138.201.140.70", connections: 184, hops: 2 },
              ].map((data, index) => (
                <li key={index} className="p-4 bg-base-100 rounded-lg border border-white border-opacity-20">
                  <div className="flex items-start">
                    <img src="/path-to-your-icon.png" alt="Country Icon" className="w-8 h-8 mr-3" />
                    <div>
                      <h3 className="text-lg font-semibold">{data.ip}</h3>
                      <ul className="menu bg-base-200 rounded-box w-full max-w-xs">
                        <li>
                          <a className="text-sm mt-2">
                            <strong>Country:</strong> {data.country}
                          </a>
                        </li>
                        <li>
                          <a className="text-sm mt-2">
                            <strong>Connections:</strong> {data.connections}
                          </a>
                        </li>
                        <li>
                          <a className="text-sm mt-2">
                            <strong>HOPS:</strong> {data.hops}
                          </a>
                        </li>
                      </ul>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>

  {/* Middle Section */}
  <div className="col-span-3 p-6 flex flex-col bg-black">
  <div className="h-full rounded-box">
    <div className="grid grid-cols-2 gap-4 h-full">
      {/* Left Section: PostCard with Scroll */}
      <div className="h-full overflow-y-auto rounded-box p-4">
        <div className="space-y-4">
          {rules.map((rule, index) => (
            <PostCard
              key={rule.id || index} // Use rule.id if it exists, otherwise use the index
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
      <div className="flex flex-col space-y-4 overflow-y-auto rounded-box p-4">
        {/* Display Monitoring Configuration */}
        <MonitoringCard
          enable={monitoring.enable}
          logFile={monitoring.log_file}
          alertThresholds={monitoring.alert_thresholds}
        />

        {/* Display Ai_Ml Configuration */}
        <AiMlCard
          modelEndpoint={AiMl.model_endpoint}
          enableAnomaly={AiMl.enable_anomaly_detection}
        />

        {/* Logging Configuration */}
        <LoggingCard
          logLevel={logging.log_level}
          logFile={logging.log_file}
          maxSize={logging.max_size}
          maxBackups={logging.max_backups}
          maxAge={logging.max_age}
        />
      </div>
    </div>
  </div>
</div>



<div className="bg-black p-4 flex flex-col items-center border-t border-white border-opacity-20 gap-4 mt-4 w-screen overflow-y-auto h-80">
  <strong className="text-white mb-4">Requests</strong>

  {/* Refresh Button */}
  <button className="btn btn-info btn-outline" onClick={handleRefresh}>
    Refresh Requests
  </button>

  {/* Loading State */}
  {loadingRequests ? (
    <div className="text-white loading loading-infinity loading-lg">Loading requests</div>
  ) : requests.length > 0 ? (
    <table className="table-auto bg-base-200 rounded-box w-full border border-white border-opacity-20 overflow-y-auto h-80">
      {/* Table Head */}
      <thead>
        <tr>
          <th className="text-white px-4 py-2">Request No</th>
          <th className="text-white px-4 py-2">Time</th>
          <th className="text-white px-4 py-2">Source</th>
          <th className="text-white px-4 py-2">Destination</th>
          <th className="text-white px-4 py-2">Details</th>
        </tr>
      </thead>

      {/* Table Body */}
      <tbody>
        {requests.slice().reverse().map((request) => (
          <tr key={request.no} className="bg-base-100 hover:bg-base-300 border-b border-white border-opacity-20">
            <td className="px-4 py-2">{request.no}</td>
            <td className="px-4 py-2">{request.time}</td>
            <td className="px-4 py-2">{request.source}</td>
            <td className="px-4 py-2">{request.destination}</td>
            <td className="px-4 py-2">
              {/* Dropdown for Additional Details */}
              <div className="dropdown dropdown-hover">
                <div tabIndex={0} role="button" className="btn btn-sm btn-outline">
                  Details
                </div>
                <ul tabIndex={0} className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow">
                  <li>
                    <a className="text-sm">
                      <strong>Protocol:</strong> {request.protocol}
                    </a>
                  </li>
                  <li>
                    <a className="text-sm">
                      <strong>Length:</strong> {request.length}
                    </a>
                  </li>
                  <li>
                    <a className="text-sm">
                      <strong>Info:</strong> {request.info}
                    </a>
                  </li>
                </ul>
              </div>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  ) : (
    <div className="text-white">No requests found.</div>
  )}
</div>

</div>
  );
};

export default DashBoard;
