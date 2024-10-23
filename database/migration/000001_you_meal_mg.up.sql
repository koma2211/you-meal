CREATE TABLE IF NOT EXISTS "clients" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "phone_number" VARCHAR(12) UNIQUE NOT NULL, 
    "name" VARCHAR(255) NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "categories" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "title" VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "meals" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "category_id" INTEGER REFERENCES categories(id) NOT NULL,
    "title" VARCHAR(255) UNIQUE NOT NULL,
    "description" TEXT NOT NULL, 
    "weight" INTEGER NOT NULL,
    "calorie" INTEGER NOT NULL,
    "price" NUMERIC(10, 2) NOT NULL,
    "image_path" VARCHAR(255) NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "ingredients" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "meal_id" INTEGER REFERENCES meals(id) NOT NULL,
    "title" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS "orders" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "client_id" INTEGER REFERENCES clients(id) NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'pending', -- pending, waiting, completed, canceled
    "total_price" NUMERIC(10, 2) NOT NULL,
    "receiving_at" timestamp with time zone NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now() 
);

CREATE TABLE IF NOT EXISTS "ordered_meals" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" INTEGER REFERENCES orders(id) NOT NULL, 
    "meal_id" INTEGER REFERENCES meals(id) NOT NULL,
    "quantity" INTEGER NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now() 
);

CREATE TABLE IF NOT EXISTS "deliveries" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" INTEGER REFERENCES orders(id) NOT NULL,
    "address" TEXT NOT NULL,
    "floor" INTEGER NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now()    
);

CREATE TABLE IF NOT EXISTS "self_pickups" (
    "id" INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    "order_id" INTEGER REFERENCES orders(id) NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT now()   
);


