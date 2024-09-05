import React, { useState, useEffect } from 'react';
import axios from 'axios';

const PostCard = ({
    application,
    allowedDomains = [],
    blockedDomains = [],
    allowedIps = [],
    blockedIps = [],
    protocols = [],
}) => {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        application,
        allowedDomains: allowedDomains.join('\n'),
        blockedDomains: blockedDomains.join('\n'),
        allowedIps: allowedIps.join('\n'),
        blockedIps: blockedIps.join('\n'),
        protocols: protocols.join('\n'),
    });
    const [error, setError] = useState('');
    const [successMessage, setSuccessMessage] = useState('');

    useEffect(() => {
        setFormData({
            application,
            allowedDomains: allowedDomains.join('\n'),
            blockedDomains: blockedDomains.join('\n'),
            allowedIps: allowedIps.join('\n'),
            blockedIps: blockedIps.join('\n'),
            protocols: protocols.join('\n'),
        });
    }, [application, allowedDomains, blockedDomains, allowedIps, blockedIps, protocols]);

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
                firewall: {
                    rules: [
                        {
                            id: 0, // Assuming ID is auto-generated or not required for POST
                            application: formData.application,
                            allowed_domains: formData.allowedDomains.split('\n').map(item => item.trim()).filter(item => item),
                            blocked_domains: formData.blockedDomains.split('\n').map(item => item.trim()).filter(item => item),
                            allowed_ips: formData.allowedIps.split('\n').map(item => item.trim()).filter(item => item),
                            blocked_ips: formData.blockedIps.split('\n').map(item => item.trim()).filter(item => item),
                            protocols: formData.protocols.split('\n').map(item => item.trim()).filter(item => item),
                        }
                    ]
                },
                monitoring: {
                    enable: true,  // This should be set according to your requirements
                    log_file: "/var/log/firewall-monitor.log",
                    alert_thresholds: {
                        abnormal_traffic: 1000,
                        blocked_attempts: 10,
                    }
                },
                ai_ml: {
                    model_endpoint: "http://localhost:5000/predict",
                    enable_anomaly_detection: true,
                },
                endpoints: [
                    // You can add endpoint details here if needed
                ],
                web_console: {
                    port: 8080,
                    allowed_ips: [
                        "192.168.1.1"
                    ],
                    blocked_ips: [],
                    admin_users: [
                        "admin"
                    ],
                },
                logging: {
                    log_level: "info",
                    log_file: "/var/log/firewall-tool.log",
                    max_size: 100,
                    max_backups: 3,
                    max_age: 30,
                },
                network: {
                    enable_deep_packet_inspection: true,
                    self_signed_tls_cert: "/path/to/self-signed-cert.pem",
                    proxy_enabled: true,
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
                
                // <div>
                //     <label className="block mb-1 text-white">Application</label>
                //     <input
                //         type="text"
                //         name="application"
                //         value={formData.application}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <label className="block mb-1 text-white">Allowed Domains</label>
                //     <textarea
                //         name="allowedDomains"
                //         value={formData.allowedDomains}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <label className="block mb-1 text-white">Blocked Domains</label>
                //     <textarea
                //         name="blockedDomains"
                //         value={formData.blockedDomains}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <label className="block mb-1 text-white">Allowed IPs</label>
                //     <textarea
                //         name="allowedIps"
                //         value={formData.allowedIps}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <label className="block mb-1 text-white">Blocked IPs</label>
                //     <textarea
                //         name="blockedIps"
                //         value={formData.blockedIps}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <label className="block mb-1 text-white">Protocols</label>
                //     <textarea
                //         name="protocols"
                //         value={formData.protocols}
                //         onChange={handleInputChange}
                //         className="p-2 mb-4 input input-bordered input-neutral w-full max-w-xs"
                //     />

                //     <button
                //         onClick={handleSave}
                //         className="rounded absolute bottom-2 right-2 btn btn-info btn-outline"
                //     >
                //         Save
                //     </button>
                //     </div>
                <div>
                <ul className="menu bg-base-200 rounded-box w-100 max-w-xs">
                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Application</strong>
                        <input
                            type="text"
                            name="application"
                            value={formData.application}
                            onChange={handleInputChange}
                            className="p-2 input input-bordered input-neutral w-100"
                        />
                        </a>
                    </li>

                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Allowed Domains</strong>
                        <textarea
                            name="allowedDomains"
                            value={formData.allowedDomains}
                            onChange={handleInputChange}
                            className="p-2 input input-bordered input-neutral w-100"
                        />
                        </a>
                    </li>

                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Blocked Domains</strong>
                        <textarea
                            name="blockedDomains"
                            value={formData.blockedDomains}
                            onChange={handleInputChange}
                            className="p-2 input input-bordered input-neutral w-100"
                        />
                        </a>
                    </li>

                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Allowed IPs</strong>
                        <textarea
                            name="allowedIps"
                            value={formData.allowedIps}
                            onChange={handleInputChange}
                            className="p-2 input input-bordered input-neutral w-100"
                        />
                        </a>
                    </li>

                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Blocked IPs</strong>
                        <textarea
                            name="blockedIps"
                            value={formData.blockedIps}
                            onChange={handleInputChange}
                            className="p-2 input input-bordered input-neutral w-100"
                        />
                        </a>
                    </li>

                    <li>
                        <a>
                        <strong className="block mb-1 text-white">Protocols</strong>
                        <textarea
                            name="protocols"
                            value={formData.protocols}
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
                            <img src="/path-to-your-icon.png" alt="Icon" className="w-8 h-8 mr-3" />
                            <div>
                                <h3 className="text-lg font-semibold">{formData.application}</h3>
                                <ul className="menu bg-base-200 rounded-box w-full max-w-xs">
                                <li>
                                    <a className="text-sm mt-2">
                                    <strong>Allowed Domains:</strong> {formData.allowedDomains}
                                    </a>
                                </li>
                                <li>
                                    <a className="text-sm mt-2">
                                    <strong>Blocked Domains:</strong> {formData.blockedDomains || 'None'}
                                    </a>
                                </li>
                                <li>
                                    <a className="text-sm mt-2">
                                    <strong>Allowed IPs:</strong> {formData.allowedIps || 'None'}
                                    </a>
                                </li>
                                <li>
                                    <a className="text-sm mt-2">
                                    <strong>Blocked IPs:</strong> {formData.blockedIps || 'None'}
                                    </a>
                                </li>
                                <li>
                                    <a className="text-sm mt-2">
                                    <strong>Protocols:</strong> {formData.protocols}
                                    </a>
                                </li>
                                </ul>
                            </div>
                            </div>
                        </li>
                    </ul>


                    
                    <button 
                    onClick={handleEditClick}
                    className="absolute bottom-2 right-2  btn btn-info btn-outline">Edit</button>
                </div>
            )}
        </div>
    );
};

export default PostCard;