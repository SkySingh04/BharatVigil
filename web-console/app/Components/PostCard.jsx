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