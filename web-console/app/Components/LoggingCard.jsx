'use client';
import React, { useState, useEffect } from 'react';
import axios from 'axios';

const LoggingCard = ({ 
    logLevel = '', 
    logFile = '', 
    maxSize = 0, 
    maxBackups = 0, 
    maxAge = 0 
}) => {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        logLevel,
        logFile,
        maxSize,
        maxBackups,
        maxAge
    });
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    useEffect(() => {
        setFormData({
            logLevel,
            logFile,
            maxSize,
            maxBackups,
            maxAge
        });
    }, [logLevel, logFile, maxSize, maxBackups, maxAge]);

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

            // Prepare the data in the correct hierarchical format
            const jsonData = {
                logging: {
                    log_level: formData.logLevel,
                    log_file: formData.logFile,
                    max_size: parseInt(formData.maxSize, 10),
                    max_backups: parseInt(formData.maxBackups, 10),
                    max_age: parseInt(formData.maxAge, 10),
                }
            };

            console.log('Sending data to server:', jsonData);

            // Send the data to the backend
            const response = await axios.post('http://localhost:8080/config', jsonData, {
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            console.log('Server Response:', response);

            if (response.status === 200) {
                setSuccessMessage('Config updated successfully!');
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
        <div className="bg-black text-white rounded-lg p-4 m-4 shadow-lg border border-white border-opacity-20 border-[1px] relative">
            {isEditing ? (
                <div>
                    <ul className="menu bg-base-200 rounded-box w-100 max-w-xs">
                        <li>
                            <a>
                                <strong className="block mb-1 text-white">Log Level</strong>
                                <input
                                    type="text"
                                    name="logLevel"
                                    value={formData.logLevel}
                                    onChange={handleInputChange}
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
                                <strong className="block mb-1 text-white">Max Size (MB)</strong>
                                <input
                                    type="number"
                                    name="maxSize"
                                    value={formData.maxSize}
                                    onChange={handleInputChange}
                                    className="p-2 input input-bordered input-neutral w-100"
                                />
                            </a>
                        </li>

                        <li>
                            <a>
                                <strong className="block mb-1 text-white">Max Backups</strong>
                                <input
                                    type="number"
                                    name="maxBackups"
                                    value={formData.maxBackups}
                                    onChange={handleInputChange}
                                    className="p-2 input input-bordered input-neutral w-100"
                                />
                            </a>
                        </li>

                        <li>
                            <a>
                                <strong className="block mb-1 text-white">Max Age (Days)</strong>
                                <input
                                    type="number"
                                    name="maxAge"
                                    value={formData.maxAge}
                                    onChange={handleInputChange}
                                    className="p-2 input input-bordered input-neutral w-100"
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
                                    <h3 className="text-lg font-semibold">Logging Configuration</h3>
                                    <ul className="menu bg-base-200 rounded-box w-full max-w-xs">
                                        <li>
                                            <a className="text-sm mt-2">
                                                <strong>Log Level:</strong> {formData.logLevel}
                                            </a>
                                        </li>
                                        <li>
                                            <a className="text-sm mt-2">
                                                <strong>Log File:</strong> {formData.logFile}
                                            </a>
                                        </li>
                                        <li>
                                            <a className="text-sm mt-2">
                                                <strong>Max Size (MB):</strong> {formData.maxSize}
                                            </a>
                                        </li>
                                        <li>
                                            <a className="text-sm mt-2">
                                                <strong>Max Backups:</strong> {formData.maxBackups}
                                            </a>
                                        </li>
                                        <li>
                                            <a className="text-sm mt-2">
                                                <strong>Max Age (Days):</strong> {formData.maxAge}
                                            </a>
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
            {error && <p className="text-red-500">{error}</p>}
            {successMessage && <p className="text-green-500">{successMessage}</p>}
        </div>
    );
};

export default LoggingCard;
