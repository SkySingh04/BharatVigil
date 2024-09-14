'use client';
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const MonitoringCard = ({ enable, logFile, alertThresholds }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
      enable,
      logFile,
      abnormalTraffic: alertThresholds?.abnormal_traffic || '',
      blockedAttempts: alertThresholds?.blocked_attempts || '',
    });
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');
  
    useEffect(() => {
      setFormData({
        enable,
        logFile,
        abnormalTraffic: alertThresholds?.abnormal_traffic || '',
        blockedAttempts: alertThresholds?.blocked_attempts || '',
      });
    }, [enable, logFile, alertThresholds]);
  
    const handleEditClick = () => {
      setIsEditing(!isEditing);
      setError('');
      setSuccessMessage('');
    };
  
    const handleInputChange = (e) => {
      const { name, value } = e.target;
      setFormData(prevState => ({
        ...prevState,
        [name]: value,
      }));
    };
  
    const handleSave = async () => {
      try {
        setError('');
        setSuccessMessage('');
  
        const jsonData = {
          monitoring: {
            enable: formData.enable,
            log_file: formData.logFile,
            alert_thresholds: {
              abnormal_traffic: Number(formData.abnormalTraffic),
              blocked_attempts: Number(formData.blockedAttempts),
            }
          }
        };
  
        const response = await axios.post('http://localhost:8080/config', jsonData, {
          headers: {
            'Content-Type': 'application/json'
          }
        });
  
        if (response.status === 200) {
          setSuccessMessage('Monitoring config updated successfully!');
          setIsEditing(false);
        } else {
          setError(response.data.error || 'Failed to update config');
        }
      } catch (error) {
        console.error('Error updating config:', error);
        setError('Error updating config');
      }
    };
  
    return (
      <div className="bg-black text-white rounded-lg w-[500px] shadow-lg border border-white border-opacity-20 border-[1px] relative">
        {isEditing ? (
          <div>
            <ul className="menu bg-base-200 rounded-box w-100 max-w-xs">
              <li>
                <a>
                  <strong className="block mb-1 text-white">Enable Monitoring</strong>
                  <input
                    type="checkbox"
                    name="enable"
                    checked={formData.enable}
                    onChange={(e) => setFormData({ ...formData, enable: e.target.checked })}
                    className="p-2 input input-bordered input-neutral w-100"
                  />
                </a>
              </li>
              <li>
                <a>
                  <strong className="block mb-1 text-white">Log File</strong>
                  <input
                    type="text"
                    name="logFile"
                    value={formData.logFile}
                    onChange={handleInputChange}
                    className="p-2 input input-bordered input-neutral w-100"
                  />
                </a>
              </li>
              <li>
                <a>
                  <strong className="block mb-1 text-white">Abnormal Traffic Threshold</strong>
                  <input
                    type="number"
                    name="abnormalTraffic"
                    value={formData.abnormalTraffic}
                    onChange={handleInputChange}
                    className="p-2 input input-bordered input-neutral w-20"
                  />
                
                  <strong className="block mb-1 text-white">Blocked Attempts Threshold</strong>
                  <input
                    type="number"
                    name="blockedAttempts"
                    value={formData.blockedAttempts}
                    onChange={handleInputChange}
                    className="p-2 input input-bordered input-neutral w-10"
                  />
                </a>
              </li>
            </ul>
            <button
              onClick={handleSave}
              className="rounded btn btn-info btn-outline absolute bottom-2 right-2"
            >
              Save
            </button>
          </div>
        ) : (
          <div>
             <ul className="menu bg-base-200 rounded-box w-100">
                        <li>
                            <div className="flex items-start">
                            <div>
             <h3 className="text-lg font-semibold">Monitoring Configuration</h3>
            <ul className="menu bg-base-200 rounded-box w-100">
              <li><a>
                <strong>Enable Monitoring:</strong> {formData.enable ? 'Yes' : 'No'}</a>
              </li>
             <li><a>
                <strong>Abnormal Traffic Threshold:</strong> {formData.abnormalTraffic}
              
                <strong>Blocked Attempts Threshold:</strong> {formData.blockedAttempts}</a>
              </li>
              <li><a>
                <strong>Log File:</strong> {formData.logFile}</a>
              </li>
            </ul>
            </div>
            </div>
            </li>
            </ul>
            <button
              onClick={handleEditClick}
              className="absolute bottom-2 right-2 btn btn-info btn-outline"
            >
              Edit
            </button>
          </div>
        )}
      </div>
    );
  };

  export default MonitoringCard;
  