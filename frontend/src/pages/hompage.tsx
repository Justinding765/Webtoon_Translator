import { useNavigate } from 'react-router-dom';
import React, { useState, useEffect } from 'react';
import { v4 as uuidv4 } from 'uuid';

import '../App.css'
const HomePage: React.FC = () => {
    const navigate = useNavigate();
    const [url, setUrl] = useState('');
    const [loading, setLoading] = useState(false);
    const [sessionId, setSessionId] = useState('');
    const [isSubmittedSuccessfully, setIsSubmittedSuccessfully] = useState(false);
    const [imageUrls, setImageUrls] = useState(null); // State to store image URLs

    useEffect(() => {
        // Check if a session ID already exists
        let currentSessionId = sessionStorage.getItem('sessionId');
        console.log(isSubmittedSuccessfully)
        // If not, generate a new one and store it in sessionStorage
        if (!currentSessionId) {
            currentSessionId = uuidv4();
            sessionStorage.setItem('sessionId', currentSessionId);
        }

        setSessionId(currentSessionId);
       
    }, [])

    useEffect(() => {
        if (isSubmittedSuccessfully && imageUrls) {
            navigate('/translation', { state: { imageUrls, isSubmittedSuccessfully } });
        }
    }, [isSubmittedSuccessfully, imageUrls, navigate]);


    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setLoading(true);  // Set loading to true when the request starts
        try {
            const response = await fetch('/api/translate_images', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url, sessionId }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            setImageUrls(data); // Store image URLs
            setIsSubmittedSuccessfully(true);
        } catch (error) {
            console.error('Error submitting URL:', error);
            setIsSubmittedSuccessfully(false); // Set to false in case of error
        } finally {
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
