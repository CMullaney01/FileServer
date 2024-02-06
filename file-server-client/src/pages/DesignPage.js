import React from 'react';
import background from "../images/background.png";
import './DesignPage.css'; // Import the CSS file

const DesignPage = () => {
    return (
        <div
            className="design-page-container" // Apply the CSS class
            style={{ backgroundImage: `url(${background})` }}
        >
            <div style={{ color: 'white' }}>Hello World</div>
        </div>
    );
}

export default DesignPage;
