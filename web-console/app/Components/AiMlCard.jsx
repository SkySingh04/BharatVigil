'use client';
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const AiMlCard = ({ modelEndpoint = '', enableAnomaly = false }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        modelEndpoint,
        enableAnomaly: enableAnomaly ? 'true' : 'false', // Ensure boolean is string
    });
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    useEffect(() => {
      setFormData({
        modelEndpoint,
        enableAnomaly: enableAnomaly ? 'true' : 'false', // Ensure boolean is string
      });
    }, [modelEndpoint, enableAnomaly]);

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
          AiMl: {
            model_endpoint: formData.modelEndpoint,
            enable_anomaly_detection: formData.enableAnomaly === 'true', // Convert string to boolean
          }
        };

        const response = await axios.post('http://localhost:8080/config', jsonData, {
          headers: {
            'Content-Type': 'application/json'
          }
        });

        if (response.status === 200) {
          setSuccessMessage('Ai_Ml config updated successfully!');
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
      <div className="bg-black text-white rounded-lg w-[450px] shadow-lg border border-white border-opacity-20 border-[1px] relative">
        {isEditing ? (
          <div>
            <ul className="menu bg-base-200 rounded-box w-100 max-w-xs">
              <li>
                <a>
                  <strong className="block mb-1 text-white">model_endpoint</strong>
                  <input
                    type="text"
                    name="model_endpoint"
                    value={formData.modelEndpoint}
                    onChange={handleInputChange}
                    className="p-2 input input-bordered input-neutral w-100"
                  />
                </a>
              </li>
              <li>
                <a>
                  <strong className="block mb-1 text-white">enable_anomaly_detection</strong>
                  <select
                    name="enable_anomaly_detection"
                    value={formData.enableAnomaly}
                    onChange={handleInputChange}
                    className="p-2 input input-bordered input-neutral w-100"
                  >
                    <option value="true">True</option>
                    <option value="false">False</option>
                  </select>
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
             <h3 className="text-lg font-semibold">AI-Ml Configuration</h3>
            <ul className="menu bg-base-200 rounded-box w-100">
              <li><a>
                <strong>model endpoint:</strong> {formData.modelEndpoint}</a>
              </li>
              <li><a>
                <strong>enable anomaly detection:</strong> {formData.enableAnomaly=== 'true' ? 'True' : 'False'}</a>
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

export default AiMlCard;
