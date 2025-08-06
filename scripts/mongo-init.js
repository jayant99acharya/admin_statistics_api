// MongoDB initialization script for Docker
// This script runs when the MongoDB container starts for the first time

// Switch to the admin_statistics database
db = db.getSiblingDB('admin_statistics');

// Create a user for the application
db.createUser({
  user: 'app_user',
  pwd: 'app_password',
  roles: [
    {
      role: 'readWrite',
      db: 'admin_statistics'
    }
  ]
});

// Create the transactions collection
db.createCollection('transactions');

// Create indexes for better performance
db.transactions.createIndex({ "createdAt": 1 });
db.transactions.createIndex({ "userId": 1 });
db.transactions.createIndex({ "userId": 1, "createdAt": 1 });
db.transactions.createIndex({ "roundId": 1 });
db.transactions.createIndex({ "type": 1 });

print('Database initialization completed successfully!');