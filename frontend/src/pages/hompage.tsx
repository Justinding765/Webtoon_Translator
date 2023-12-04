// src/pages/HomePage.tsx
import React, { useState } from 'react';
import '../App.css';

const HomePage: React.FC = () => {
    const [url, setUrl] = useState('');

    const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        try {
            const response = await fetch('http://localhost:5000/web_scrapper', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ url: url }), // Sending the URL as JSON
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            console.log('Response from server:', data);
            // Handle the response data
        } catch (error) {
            console.error('Error submitting URL:', error);
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
                    <button type="submit" className="homepage-button">Submit</button>
                </form>
            </div>
        </div>
    );
};

export default HomePage;
