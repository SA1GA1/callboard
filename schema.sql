CREATE TABLE listings (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    seller_id BIGINT NOT NULL,
    price INTEGER NOT NULL,
    place TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_sold BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE specifications (
    listing_id BIGINT PRIMARY KEY,
    category TEXT,
    brand TEXT,
    model TEXT,
    condition TEXT,
    selling_reason TEXT,
    CONSTRAINT fk_listings_specifications
        FOREIGN KEY (listing_id)
        REFERENCES listings (id)
        ON DELETE CASCADE
);

CREATE TABLE photos (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    listing_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    CONSTRAINT fk_listings_photos
        FOREIGN KEY (listing_id)
        REFERENCES listings (id)
        ON DELETE CASCADE
);