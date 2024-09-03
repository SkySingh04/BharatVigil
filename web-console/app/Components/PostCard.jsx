import React, { useState } from 'react';

const PostCard = ({ 
    application, 
    allowedDomains, 
    blockedDomains, 
    allowedIps, 
    blockedIps, 
    protocols 
}) => {
    const [isEditing, setIsEditing] = useState(false);
    const [formData, setFormData] = useState({
        application,
        allowedDomains: allowedDomains.join(', '),
        blockedDomains: blockedDomains.join(', '),
        allowedIps: allowedIps.join(', '),
        blockedIps: blockedIps.join(', '),
        protocols: protocols.join(', '),
    });

    const handleEditClick = () => {
        setIsEditing(!isEditing);
    };

    const handleInputChange = (e) => {
        const { name, value } = e.target;
        setFormData({
            ...formData,
            [name]: value,
        });
    };

    const handleSave = () => {
        // Here you could send the updated data to your backend or simply update the UI
        setIsEditing(false);
        console.log('Updated data:', formData);
    };

    return (
        <div className="bg-black text-white rounded-lg p-4 m-4 shadow-lg border border-white border-opacity-20 border-[1px] relative">
            {isEditing ? (
                <div>
                    <input
                        type="text"
                        name="application"
                        value={formData.application}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded bg-custom-gray"
                    />
                    <textarea
                        name="allowedDomains"
                        value={formData.allowedDomains}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded rounded bg-custom-gray"
                    />
                    <textarea
                        name="blockedDomains"
                        value={formData.blockedDomains}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded rounded bg-custom-gray"
                    />
                    <textarea
                        name="allowedIps"
                        value={formData.allowedIps}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded rounded bg-custom-gray"
                    />
                    <textarea
                        name="blockedIps"
                        value={formData.blockedIps}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded rounded bg-custom-gray"
                    />
                    <textarea
                        name="protocols"
                        value={formData.protocols}
                        onChange={handleInputChange}
                        className="w-full mb-2 p-2 text-black rounded rounded bg-custom-gray"
                    />
                    <button
                        onClick={handleSave}
                        className="bg-blue-500 text-white p-2 rounded rounded bg-green-500"
                    >
                        Save
                    </button>
                </div>
            ) : (
                <div>
                    <div className="flex items-start">
                        <img src="/path-to-your-icon.png" alt="Icon" className="w-8 h-8 mr-3" />
                        <div>
                            <h3 className="text-lg font-semibold">{formData.application}</h3>
                            <p className="text-sm mt-2">
                                <strong>Allowed Domains:</strong> {formData.allowedDomains}
                            </p>
                            <p className="text-sm mt-2">
                                <strong>Blocked Domains:</strong> {formData.blockedDomains || 'None'}
                            </p>
                            <p className="text-sm mt-2">
                                <strong>Allowed IPs:</strong> {formData.allowedIps || 'None'}
                            </p>
                            <p className="text-sm mt-2">
                                <strong>Blocked IPs:</strong> {formData.blockedIps || 'None'}
                            </p>
                            <p className="text-sm mt-2">
                                <strong>Protocols:</strong> {formData.protocols}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={handleEditClick}
                        className="absolute bottom-2 right-2 bg-blue-500 text-white p-2 rounded "
                    >
                        Edit
                    </button>
                </div>
            )}
        </div>
    );
};
export default PostCard;
