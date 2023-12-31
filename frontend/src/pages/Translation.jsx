import '../App.css';
import React, { useState } from 'react';
import { useLocation } from "react-router-dom";



const MyComponent = () => {
  // State to toggle navbar visibility
  const [showNavbar, setShowNavbar] = useState(false);
  const location = useLocation();
  const { imageUrls } = location.state || { imageUrls: [] };
  const [loading, setLoading] = useState(false);  // New state for loading
  console.log(imageUrls)
  const handleDownloadPdf = async () => {
    setLoading(true);  // Set loading to true when the request s
    const requestUrl = 'http://localhost:5000/download_pdf';

    try {
        const response = await fetch(requestUrl);
        console.log(response)
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
