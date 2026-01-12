CREATE TABLE IF NOT EXISTS books (
                                     id SERIAL PRIMARY KEY,
                                     title TEXT NOT NULL,
                                     author TEXT NOT NULL,
                                     year INT,
                                     isbn TEXT UNIQUE,
                                     out_of_stock BOOLEAN DEFAULT FALSE,
                                     rating INT,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
