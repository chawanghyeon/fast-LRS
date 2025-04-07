const express = require('express');
const { v4: uuidv4 } = require('uuid');
const { Pool } = require('pg');

const app = express();
app.use(express.json());

const pool = new Pool({
    host: 'postgres',
    user: 'benchmark',
    password: 'benchmark',
    database: 'benchmark',
    port: 5432
});

function validateUser(user) {
    const emailPattern = /^[^\s@]+@example\.com$/;
    const bioRegex = /(engineer|developer|programmer)/i;

    if (!user.name || user.name.length < 2 || user.name.length > 50) {
        return 'Name length must be between 2 and 50';
    }
    if (!emailPattern.test(user.email)) {
        return 'Email must be @example.com';
    }
    if (typeof user.age !== 'number' || user.age < 0 || user.age > 120) {
        return 'Invalid age';
    }
    if (!user.bio || user.bio.length < 10 || !bioRegex.test(user.bio)) {
        return 'Invalid bio';
    }
    if (!Array.isArray(user.interests) || user.interests.length === 0) {
        return 'At least one interest required';
    }
    return null;
}

app.post('/data', async (req, res) => {
    const user = req.body;
    const error = validateUser(user);
    if (error) return res.status(400).json({ error });

    const id = uuidv4();

    try {
        await pool.query(
            'INSERT INTO users (id, name, email) VALUES ($1, $2, $3)',
            [id, user.name, user.email]
        );
        res.json({ id });
    } catch (err) {
        console.error(err);
        res.status(500).json({ error: 'Database error' });
    }
});

app.listen(8000, () => {
    console.log('Node.js server running on port 8000');
});
