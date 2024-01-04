import { useNavigate } from 'react-router-dom';
import React, { useState } from 'react';
import '../App.css'
const HomePage: React.FC = () => {
    const navigate = useNavigate();
    const [url, setUrl] = useState('');
    const [loading, setLoading] = useState(false);  // New state for loading
    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLoading(true);  // Set loading to true when the request starts
        console.log("/api/d")
        try {
            const response = await fetch('http://localhost:5000/translate_images',  {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url }),
            });


            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            setLoading(false);  // Set loading to false when the request completes
            const data = await response.json();
            
            // Redirect and pass the image URLs
            navigate('/translation', { state: { imageUrls: data } });
        } catch (error) {
            console.error('Error submitting URL:', error);
            setLoading(false);
        }
    };


    return (
        <div>
            <nav className="navbar">
                <h1>Webtoon Translator</h1>
            </nav>
            <div className="homepage-container">
                <form onSubmit={handleSubmit} className="homepage-form">
                    <label htmlFor="url-input" className="homepage-label">Enter Webtoon URL:</label>
                    <input
                        type="text"
                        id="url-input"
                        value={url}
                        onChange={(e) => setUrl(e.target.value)}
                        className="homepage-input"
                        placeholder="https://comic.naver.com/webtoon/detail?titleId=641253&no=472&week=fri"
                    />
                     <button type="submit" className="homepage-button" disabled={loading}>
                        {loading ? <div className="loading-spinner"></div> : "Submit"}
                    </button>
                </form>
            </div>
        </div>
    );
};

export default HomePage;
