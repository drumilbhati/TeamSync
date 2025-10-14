import express from 'express';
import connectDB from './db/db_connect.js';
import dotenv from 'dotenv';

dotenv.config();

const app = express();

const port = process.env.PORT || 3000;

connectDB();

app.listen(port, () => {
	console.log(`Server listening on port ${port}`);
});