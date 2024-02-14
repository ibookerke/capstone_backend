CREATE SCHEMA IF NOT EXISTS "product";

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "product"."property_types" (
    "id" bigserial PRIMARY KEY,
    "type" text UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "product"."property_subtypes" (
    "id" bigserial PRIMARY KEY,
    "property_type_id" bigint NOT NULL,
    "subtype" text NOT NULL,
    FOREIGN KEY ("property_type_id") REFERENCES "product"."property_types" ("id")
);

CREATE TABLE IF NOT EXISTS "product"."properties" (
    "id" bigserial PRIMARY KEY,
    "name" text NOT NULL,
    "description" text NOT NULL,
    "area" float8 NOT NULL,
    "baths" int NOT NULL,
    "beds" int NOT NULL,
    "rooms" int NOT NULL,
    "location" text NOT NULL,
    "lat_lang" point NOT NULL,
    "furnishing_status" int NOT NULL,
    "property_subtype_id" bigint NOT NULL,
    FOREIGN KEY ("property_subtype_id") REFERENCES "product"."property_subtypes" ("id")
);

CREATE TABLE IF NOT EXISTS "product"."amenities" (
    "id" bigserial PRIMARY KEY,
    "name" text UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS "product"."property_amenities" (
    "id" bigserial PRIMARY KEY,
    "property_id" bigint NOT NULL,
    "amenity_id" bigint NOT NULL,
    FOREIGN KEY ("property_id") REFERENCES "product"."properties" ("id"),
    FOREIGN KEY ("amenity_id") REFERENCES "product"."amenities" ("id")
);

CREATE TABLE IF NOT EXISTS "product"."property_images" (
    "id" bigserial PRIMARY KEY,
    "property_id" bigint NOT NULL,
    "image" text NOT NULL,
    "is_main" bool NOT NULL DEFAULT false,
    FOREIGN KEY ("property_id") REFERENCES "product"."properties" ("id")
);

-- TODO: change is_active to status?
CREATE TABLE IF NOT EXISTS "product"."investment_opportunities" (
    "id" bigserial PRIMARY KEY,
    "property_id" bigint NOT NULL,
    "type" int NOT NULL, -- example: 1 means hold, 2 means flip
    "shares_amount" bigint NOT NULL, -- example: 100000
    "price_per_share" bigint NOT NULL, -- example: 10.00
    "is_active" bool NOT NULL DEFAULT FALSE, -- if io available for investing
    "currency" text NOT NULL, -- currency of the investment
    FOREIGN KEY ("property_id") REFERENCES "product"."properties" ("id")
);

CREATE TABLE IF NOT EXISTS "product"."outbox" (
    "id" bigserial PRIMARY KEY,
    "event_id" uuid NOT NULL,
    "topic_name" text NOT NULL,
    "payload" bytea NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "sent_at" timestamptz DEFAULT NULL,
    CONSTRAINT "outbox_topic_name_event_id_unique" UNIQUE ("topic_name", "event_id")
);

CREATE TABLE IF NOT EXISTS "product"."inbox" (
    "id" bigserial PRIMARY KEY,
    "event_id" uuid NOT NULL,
    "topic_name" text NOT NULL,
    "payload" bytea NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT "inbox_topic_name_event_id_unique" UNIQUE ("topic_name", "event_id")
);

-- TODO: add index on investment_opportunity_id
CREATE TABLE IF NOT EXISTS "product"."financials" (
    "id" bigserial PRIMARY KEY,
    "investment_opportunity_id" bigint NOT NULL,
    "label" varchar NOT NULL,
    "items" jsonb NOT NULL,
    "type" varchar,
    FOREIGN KEY ("investment_opportunity_id") REFERENCES "product"."investment_opportunities" ("id"),
    CONSTRAINT "financials_label_investment_opportunity_id_unique" UNIQUE("investment_opportunity_id", "label")
);

CREATE TABLE IF NOT EXISTS "product"."user_investments" (
    "id" bigserial PRIMARY KEY,
    "user_id" varchar NOT NULL,
    "investment_opportunity_id" bigint NOT NULL,
    "shares_amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    FOREIGN KEY ("investment_opportunity_id") REFERENCES "product"."investment_opportunities" ("id"),
    CONSTRAINT "user_investments_user_id_investment_opportunity_id_unique" UNIQUE("user_id", "investment_opportunity_id")
);

-- Seeding
INSERT INTO product.property_types(type) VALUES
    ('Residential'),
    ('Commercial')
;

INSERT INTO product.property_subtypes(property_type_id, subtype) VALUES
    ((select id from product.property_types where type = 'Residential'), 'Villa'),
    ((select id from product.property_types where type = 'Residential'), 'Apartments'),
    ((select id from product.property_types where type = 'Residential'), 'Villa'),
    ((select id from product.property_types where type = 'Residential'), 'Townhouse'),
    ((select id from product.property_types where type = 'Residential'), 'Penthouse'),
    ((select id from product.property_types where type = 'Residential'), 'Villa Compaund'),
    ((select id from product.property_types where type = 'Residential'), 'Hotel Apartments'),
    ((select id from product.property_types where type = 'Residential'), 'Residential Plot'),
    ((select id from product.property_types where type = 'Residential'), 'Residential Floor'),
    ((select id from product.property_types where type = 'Residential'), 'Residential Building'),
    ((select id from product.property_types where type = 'Commercial'), 'Office'),
    ((select id from product.property_types where type = 'Commercial'), 'Shop'),
    ((select id from product.property_types where type = 'Commercial'), 'Warehouse'),
    ((select id from product.property_types where type = 'Commercial'), 'Labour Camp'),
    ((select id from product.property_types where type = 'Commercial'), 'Commercial Villa'),
    ((select id from product.property_types where type = 'Commercial'), 'Bulk Unit'),
    ((select id from product.property_types where type = 'Commercial'), 'Commercial Plot'),
    ((select id from product.property_types where type = 'Commercial'), 'Commercial Floor'),
    ((select id from product.property_types where type = 'Commercial'), 'Commercial Building'),
    ((select id from product.property_types where type = 'Commercial'), 'Showroom'),
    ((select id from product.property_types where type = 'Commercial'), 'Other Commercial')
;

INSERT INTO product.amenities(name) VALUES
    ('Gym'),
    ('Pool'),
    ('Parking'),
    ('Security System'),
    ('Green Spaces or Parks'),
    ('Playground'),
    ('Pet-Friendly Features'),
    ('Balcony or Patio'),
    ('Elevator'),
    ('Laundry Facilities'),
    ('Air Conditioning and Heating'),
    ('Storage Units'),
    ('High-Speed Internet and Cable'),
    ('Concierge Services'),
    ('BBQ or Picnic Area'),
    ('Tennis or Basketball Courts'),
    ('Yoga or Meditation Studio'),
    ('On-Site Maintenance'),
    ('Electric Vehicle (EV) Charging Stations'),
    ('Wheelchair Accessibility')
;