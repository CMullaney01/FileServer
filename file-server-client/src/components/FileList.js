import React, { useState, useEffect } from "react";
import axios from "axios";
import "./FileList.css"; // Import the CSS file for styling

const FileList = () => {
  const [files, setFiles] = useState([]);

  useEffect(() => {
    const fetchFiles = async () => {
      try {
        // Use relative URL for both HTTP and HTTPS
        const response = await axios.get("http://localhost:8080/list", {
          withCredentials: true,
        });
    
        if (response.ok) {
          // Assuming the response is an array of file names
          setFiles(response.data);
        } else if (response.status === 401) {
          console.error('Unauthorized access. Redirecting to login...');
          // Redirect to login or handle unauthorized access
        } else {
          // Handle other error cases
          console.error('Error fetching files:', response.statusText);
        }
      } catch (error) {
        console.error('Error fetching files:', error.message);
      }
    };

    fetchFiles();
  }, []);

  const handleDownload = (fileName) => {
    // Construct the download URL
    const downloadURL = `http://localhost:8080/download?filename=${fileName}`;

    // Create a hidden link and trigger the download
    const link = document.createElement("a");
    link.href = downloadURL;
    link.download = fileName;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  return (
    <div className="file-list-container">
      <h2>File List</h2>
      <ul className="file-list">
        {files.map((file, index) => (
          <li key={index} onClick={() => handleDownload(file)}>
            {file}
          </li>
        ))}
      </ul>
    </div>
  );
};

export default FileList;
