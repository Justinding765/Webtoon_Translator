import '../App.css';
import React, { useState } from 'react';
import { useLocation } from "react-router-dom";


// Function to import all images in a directory
function importAll(r) {
  // Map and sort based on the numerical part of the file name
  return r.keys()
    .map(r)
    .sort((a, b) => {
      const matchA = a.match(/\d+/); // Extract number from string
      const matchB = b.match(/\d+/);
      const numA = matchA ? parseInt(matchA[0], 10) : 0;
      const numB = matchB ? parseInt(matchB[0], 10) : 0;
      return numA - numB; // Sort numerically
    });
}


// Import images from a specific directory
const images = importAll(require.context('../../src/Static/Images', false, /\.(png|jpe?g|svg)$/));
  
const MyComponent = () => {
  // State to toggle navbar visibility
  const [showNavbar, setShowNavbar] = useState(false);
  const location = useLocation();
  const { imageUrls } = location.state || { imageUrls: [] };

  const handleDownloadHtml = async () => {
    // Manually create the HTML content
    let htmlContent = `
        <html>
            <body>`
    imageUrls.forEach((item, index) => {
      let url = item['URL']
      htmlContent += `<div><img src="${url}" alt="Image ${index}" /></div>\n`;
 
    });
    console.log(htmlContent)
    
    htmlContent+= `</body>
        </html>`
        

    // Create a Blob from the HTML string and download it
    const blob = new Blob([htmlContent], { type: 'text/html' });
    const downloadLink = document.createElement('a');
    downloadLink.href = URL.createObjectURL(blob);
    downloadLink.download = 'output.html';
    document.body.appendChild(downloadLink);
    downloadLink.click();
    document.body.removeChild(downloadLink);
};


  return (
    <div onMouseOver={() => setShowNavbar(true)} onMouseOut={() => setShowNavbar(false)}>
      <div className={`navbar ${showNavbar ? '' : 'navbar-hidden'}`}>
        <button onClick={handleDownloadHtml} className="download-button">Download HTML</button>
      </div>
      <div className="content">
        {images.map((img, index) => (
          <div className="image-container" key={index}>
            <img src={img} alt={`Translated image ${index}`} />
          </div>
        ))}
      </div>
    </div>
  );
};

export default MyComponent;
