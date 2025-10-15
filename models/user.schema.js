import mongoose from 'mongoose';

const userSchema = new mongoose.Schema({
	email: {
		type: String,
		require: true,
	}, 
	password: {
		type: String,
		require: true,
	},
	name: {
		type: String,
		require: true,
	},
	role: {
		type: String,
		enum: ['admin', 'manager', 'member'],
		default: 'member',
	},
	is_verified: {
		type: Boolean,
		default: false
	},
	reset_token: {
		type: String,
	},
}, {timestamps: true}
);

export default mongoose.model('User', userSchema);