"use client";

import { useState, useRef, useEffect } from "react";

export default function Terminal() {
  const [height, setHeight] = useState(200); // Initial terminal height
  const [isDragging, setIsDragging] = useState(false); // Track dragging state
  const [loadingRequests, setLoadingRequests] = useState(true); // Manage loading state
  const [requests, setRequests] = useState<any[]>([]); // Manage requests list

  const terminalRef = useRef<HTMLDivElement>(null);
  const dragHandleRef = useRef<HTMLDivElement>(null);

  // Fetch requests from the server
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

  // SSE connection to get real-time updates from the server
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

  // Handle mouse drag to resize terminal upwards
  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (isDragging && terminalRef.current) {
        const terminalBottom = terminalRef.current.getBoundingClientRect().bottom;
        const newHeight = terminalBottom - e.clientY; // Calculate new height based on bottom drag
        setHeight(Math.max(100, Math.min(newHeight, window.innerHeight - 50))); // Restrict min/max heights
      }
    };

    const handleMouseUp = () => {
      setIsDragging(false);
    };

    if (isDragging) {
      document.addEventListener("mousemove", handleMouseMove);
      document.addEventListener("mouseup", handleMouseUp);
    }

    return () => {
      document.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("mouseup", handleMouseUp);
    };
  }, [isDragging]);

  // Initiate drag action
  const handleMouseDown = () => {
    setIsDragging(true);
  };

  return (
    <div
      ref={terminalRef}
      className="bg-black text-white font-mono text-sm absolute bottom-0 w-full"
      style={{ height: `${height}px` }}
    >
      <div
        ref={dragHandleRef}
        className="h-2 bg-gray-700 cursor-ns-resize"
        onMouseDown={handleMouseDown}
      />
      <div className="p-4 overflow-auto" style={{ height: `${height - 8}px` }}>
        {loadingRequests ? (
          <div className="text-white loading loading-infinity loading-lg">
            Loading requests...
          </div>
        ) : requests.length > 0 ? (
          <div className="overflow-y-auto w-full h-full">
            {requests.slice().reverse().map((request) => (
              <div
                key={request.no}
                className="bg-base-100 hover:bg-base-300 border-b border-white border-opacity-20 text-white p-4"
              >
                <div className="flex flex-wrap items-center gap-4">
                  <span className="text-sm">
                    <strong>Time:</strong> {request.time}
                  </span>
                  <span className="text-sm">
                    <strong>Source:</strong> {request.source} {"-->"}
                    <strong> Destination:</strong> {request.destination}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-white">No requests found.</div>
        )}
      </div>
    </div>
  );
}
