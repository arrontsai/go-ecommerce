// Create collections for our services
db = db.getSiblingDB('ecommerce');

// Create users collection for auth service
db.createCollection('users');
db.users.createIndex({ "email": 1 }, { unique: true });

// Create products collection for product service
db.createCollection('products');
db.products.createIndex({ "sku": 1 }, { unique: true });

// Create categories collection for product service
db.createCollection('categories');
db.categories.createIndex({ "name": 1 }, { unique: true });

// Create carts collection for cart service
db.createCollection('carts');
db.carts.createIndex({ "userId": 1 }, { unique: true });

// Insert some sample data
db.categories.insertMany([
  { name: "Electronics", description: "Electronic devices and accessories" },
  { name: "Clothing", description: "Apparel and fashion items" },
  { name: "Books", description: "Books, e-books, and publications" },
  { name: "Home & Kitchen", description: "Home appliances and kitchen items" }
]);

db.products.insertMany([
  {
    name: "Smartphone X",
    description: "Latest smartphone with advanced features",
    price: 999.99,
    sku: "PHONE-X-001",
    categoryId: db.categories.findOne({ name: "Electronics" })._id,
    inventory: 100,
    images: ["phone-x-1.jpg", "phone-x-2.jpg"],
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Laptop Pro",
    description: "High-performance laptop for professionals",
    price: 1499.99,
    sku: "LAPTOP-PRO-001",
    categoryId: db.categories.findOne({ name: "Electronics" })._id,
    inventory: 50,
    images: ["laptop-pro-1.jpg", "laptop-pro-2.jpg"],
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Cotton T-Shirt",
    description: "Comfortable cotton t-shirt",
    price: 19.99,
    sku: "TSHIRT-001",
    categoryId: db.categories.findOne({ name: "Clothing" })._id,
    inventory: 200,
    images: ["tshirt-1.jpg", "tshirt-2.jpg"],
    createdAt: new Date(),
    updatedAt: new Date()
  },
  {
    name: "Programming in Go",
    description: "Learn Go programming language",
    price: 39.99,
    sku: "BOOK-GO-001",
    categoryId: db.categories.findOne({ name: "Books" })._id,
    inventory: 75,
    images: ["go-book-1.jpg"],
    createdAt: new Date(),
    updatedAt: new Date()
  }
]);
