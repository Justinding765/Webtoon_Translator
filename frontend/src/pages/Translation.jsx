import '../App.css';
import React, { useState, useEffect } from 'react';
import { useLocation } from "react-router-dom";
import { useNavigate } from 'react-router-dom';


const MyComponent = () => {
  // State to toggle navbar visibility
  const [showNavbar, setShowNavbar] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { imageUrls, isSubmittedSuccessfully } = location.state || {};
  const [loading, setLoading] = useState(false);  // New state for loading

  useEffect(() => {
    if (!isSubmittedSuccessfully) {
        navigate('/'); // Replace with the route to redirect to
    }
}, [isSubmittedSuccessfully, navigate]);


useEffect(() => {
  let currentSessionId = sessionStorage.getItem('sessionId');
  let inactivityTimeout;

  const handleInactivity = () => {
    // API call on inactivity
    navigator.sendBeacon('/api/session-inactive', JSON.stringify({ sessionId: currentSessionId }));
    // Navigate to the root route after the API call
    navigate('/');
  };

  const resetInactivityTimeout = () => {
    // Clear the existing timeout
    clearTimeout(inactivityTimeout);
    // Set a new timeout
    inactivityTimeout = setTimeout(handleInactivity, 5 * 60 * 1000); // 5 minutes
  };

  // Reset the inactivity timer on these events
  window.addEventListener('mousemove', resetInactivityTimeout);
  window.addEventListener('keydown', resetInactivityTimeout);
  window.addEventListener('scroll', resetInactivityTimeout);
  window.addEventListener('click', resetInactivityTimeout);

  // Set the initial timeout
  resetInactivityTimeout();

  return () => {
    clearTimeout(inactivityTimeout);
    window.removeEventListener('mousemove', resetInactivityTimeout);
    window.removeEventListener('keydown', resetInactivityTimeout);
    window.removeEventListener('scroll', resetInactivityTimeout);
    window.removeEventListener('click', resetInactivityTimeout);
  };
}, [navigate]); // include navigate in the dependency array





  const handleDownloadPdf = async () => {
    setLoading(true);  // Set loading to true when the request s
    const requestUrl = '/api/download_pdf';
    try {
        const response = await fetch(requestUrl);
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        const blob = await response.blob();
        const downloadUrl = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = downloadUrl;
        link.download = 'output.pdf';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        setLoading(false);

    } catch (error) {
        console.error('Error downloading PDF:', error);
        setLoading(false);

    }
    
};
  if (!isSubmittedSuccessfully) {
    // Optionally, you can render something else or redirect
    return <div>Not authorized to view this page</div>;
  }

  return (
    <div onMouseOver={() => setShowNavbar(true)} onMouseOut={() => setShowNavbar(false)}>
      <div className={`navbar ${showNavbar ? '' : 'navbar-hidden'}`}>
        <button onClick={handleDownloadPdf} className="download-button" disabled={loading}>
          {loading ? <div className="loading-spinner"></div> : "Download PDF"}
        </button>
                        
      </div>
      <div className="content">
        {imageUrls.map((img, index) => (
          <div className="image-container" key={index}>
            <img src={img} alt={`Translated image ${index}`} />
          </div>
        ))}
      </div>
    </div>
  );
};

export default MyComponent;
