import React from 'react';

const PostCard = ({ application, allowedDomains, blockedDomains, allowedIps, blockedIps, protocols }) => {
    return (
      <div className="bg-black text-white rounded-lg p-4 m-4 shadow-lg border border-white border-opacity-20 border-[1px]">
        <div className="flex items-start">
          <img src="/path-to-your-icon.png" alt="Icon" className="w-8 h-8 mr-3" />
          <div>
            <h3 className="text-lg font-semibold">{application}</h3>
            <p className="text-sm mt-2">
              <strong>Allowed Domains:</strong> {allowedDomains.join(', ')}
            </p>
            <p className="text-sm mt-2">
              <strong>Blocked Domains:</strong> {blockedDomains.length > 0 ? blockedDomains.join(', ') : 'None'}
            </p>
            <p className="text-sm mt-2">
              <strong>Allowed IPs:</strong> {allowedIps.length > 0 ? allowedIps.join(', ') : 'None'}
            </p>
            <p className="text-sm mt-2">
              <strong>Blocked IPs:</strong> {blockedIps.length > 0 ? blockedIps.join(', ') : 'None'}
            </p>
            <p className="text-sm mt-2">
              <strong>Protocols:</strong> {protocols.join(', ')}
            </p>
          </div>
        </div>
      </div>
    );
}

export default PostCard;
