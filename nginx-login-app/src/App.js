import React from 'react';  
import Register from './Register';  
import Login from './Login';  
import './App.css'; 

const App = () => {  
  return (  
    <div>  
      <h1 className="neon">EmbedUR-Internal Tools and Integration Team</h1>  
      <h2 className="neon">Assessment Task</h2>  
      <Register />  
      <Login />  
    </div>  
  );  
};  

export default App;