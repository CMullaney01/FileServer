import React, { useState, useEffect } from "react";
import axios from "axios";

const FileList = () => {
  const [files, setFiles] = useState([]);

  useEffect(() => {
    const fetchFiles = async () => {
      try {
        // Use relative URL for both HTTP and HTTPS
        const response = await axios.get("http://localhost:8080/list");
        console.log("Response from server:", response);

        // Assuming the response is an array of file names
        setFiles(response.data);
      } catch (error) {
        console.error("Error fetching files:", error);
      }
    };

    fetchFiles();
  }, []);

  return (
    <div>
      <h2>File List</h2>
      <ul>
        {files.map((file, index) => (
          <li key={index}>{file}</li>
        ))}
      </ul>
    </div>
  );
};

export default FileList;
