"use client" ;
import React, { useEffect, useState } from 'react';
import PostCard from './PostCard';
import yaml from 'js-yaml';

const DashBoard = () => {
  const [rules, setRules] = useState([]);

  useEffect(() => {
    // Fetch data from backend
    const fetchConfig = async () => {
      try {
        const response = await fetch('http://localhost:8080/config'); // Adjust URL based on your backend setup
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const yamlText = await response.text();
        const data = yaml.load(yamlText);

        // Accessing firewall rules from the fetched data
        const firewallRules = data.firewall?.rules || [];
        setRules(firewallRules);
      } catch (error) {
        console.error('Error fetching config:', error);
      }
    };

    fetchConfig();
  }, []);



  return (
    <div className="min-h-screen bg-black grid grid-cols-3 gap-4">
      
      {/* Left Section */}
      <div className="col-span-1 bg-black p-4 flex flex-col items-center border-r border-white">
        {/* Your existing left section code */}
      </div>

      {/* Middle Section */}
      <div className="col-span-1 p-6 flex flex-col items-center bg-black">
        <div className="w-full space-y-4">
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

      {/* Right Section */}
      <div className="col-span-1 bg-black p-4 flex flex-col items-center border-l border-white">
        {/* Your existing right section code */}
      </div>
    </div>
  );
};

export default DashBoard;
