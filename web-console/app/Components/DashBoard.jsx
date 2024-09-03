"use client" ;
import React, { useEffect, useState } from 'react';
import PostCard from './PostCard';
import yaml from 'js-yaml';
import { FaBell, FaExclamationTriangle, FaCog, FaSignOutAlt } from 'react-icons/fa';


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

    
      <div className="col-span-1 border-r border-white">
        {/* Your existing left section code */}

            <div className="flex h-screen bg-black text-white">
              {/* Left Navbar with Icons */}
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

              {/* Sidebar */}
              <div className="bg-black p-8">
                <div className="flex items-center justify-between mb-4">
                  
                  <img src="/secure.jpeg" className='w-20 h-20 mb-2 rounded-full flex items-center' ></img>
                  <div className="flex items-center">
                    <div>
                    {/* <div className="p-4">
                        {/* Example usage of the SecurityStatus component */
                        /* <SecurityStatus />
                      </div> */} 
                      
                      <button className="btn btn-lg btn-success">Secure</button>
                    </div>
                  </div>
                </div>
                <div className="mb-4">
                  <label className="input input-bordered flex items-center gap-2">
                    <input type="text" className="grow" placeholder="Search" />
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      viewBox="0 0 16 16"
                      fill="currentColor"
                      className="h-4 w-4 opacity-70">
                      <path
                        fillRule="evenodd"
                        d="M9.965 11.026a5 5 0 1 1 1.06-1.06l2.755 2.754a.75.75 0 1 1-1.06 1.06l-2.755-2.754ZM10.5 7a3.5 3.5 0 1 1-7 0 3.5 3.5 0 0 1 7 0Z"
                        clipRule="evenodd" />
                    </svg>
                  </label> 
                </div>
                <div>
                <ul className="menu bg-base-200 rounded-box w-full">
                  {['Microsoft Edge', 'Portmaster Core Service', 'System DNS Client', 'Firefox'].map((app, index) => (
                    <li key={index}>
                      <a className="flex justify-between items-center p-4 bg-black rounded mb-2">
                        <span>{app}</span>
                        <div
                          className="radial-progress bg-primary text-primary-content border-primary border-2"
                          style={{ "--value": 50 }}
                          role="progressbar">
                          50%
                        </div>
                      </a>
                    </li>
                  ))}
                </ul>

              </div>
              {/* Connections Information */}
              <div className="mt-4 grid grid-cols-1 gap-4">
                  {[
                    { country: 'Canada', ip: '78.138.17.182', connections: 2, hops: 2 },
                    { country: 'USA', ip: '5.34.178.198', connections: 45, hops: 2 },
                    { country: 'France', ip: '141.95.158.73', connections: 6, hops: 2 },
                    { country: 'Germany', ip: '138.201.140.70', connections: 184, hops: 2 },
                  ].map((data, index) => (
                    <div key={index} className="bg-gray-900 p-4 border border-white border-opacity-20 border-[1px] relative rounded">
                      <div className="text-sm">{data.country}</div>
                      <div className="text-lg font-semibold">{data.ip}</div>
                      <div className="text-sm">
                        {data.connections} Connections, HOPS: {data.hops}
                      </div>
                    </div>
                  ))}
                </div>
                </div>

             
            </div>

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
