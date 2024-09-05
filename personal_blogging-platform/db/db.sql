
CREATE TABLE IF NOT EXISTS users(
    ID uuid DEFAULT gen_random_uuid() PRIMARY KEY,     
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    user_password VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS article(
    ID uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    article_title VARCHAR(200) NOT NULL,
    article_content VARCHAR(20000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at DATE ,
    user_id uuid ,
    FOREIGN KEY (user_id) REFERENCES users (ID)
)

