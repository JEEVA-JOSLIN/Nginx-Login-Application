import React, { useState } from 'react';
import axios from 'axios';
import './Register.css'; 

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await axios.post('/register', { username, password });
            alert('Registration successful!');
        } catch (error) {
            alert('Error registering: ' + error.response.data);
        }
    };

    return (
        <div className="container">
            <form onSubmit={handleSubmit} className="form">
                <h2 className="heading">REGISTER</h2>
                <input
                    type="text"
                    placeholder="Username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                    className="input"
                />
                <input
                    type="password"
                    placeholder="Password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    className="input"
                />
                <button type="submit" className="button">Register</button>
            </form>
        </div>
    );
};

export default Register;