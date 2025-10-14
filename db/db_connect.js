import dotenv from 'dotenv';
import mongoose from 'mongoose';

dotenv.config();

const connectDB = async () => {
	try {
		const connect = await mongoose.connect(process.env.MONGO_URI);
		console.log('MongoDB connected');
	} catch (error) {
		console.error(error.message);
		process.exit(1);
	}
}

export default connectDB;