-- Trigger function to update updated_at timestamps
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- user
CREATE TABLE "user" (
  "id" uuid NOT NULL,
  "username" varchar(255) NOT NULL UNIQUE,
  "name" varchar(255) NOT NULL,
  "bio" text,
  "password" varchar(255) NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  "deleted_at" timestamp with time zone,
  PRIMARY KEY ("id")
);
CREATE INDEX "user_index_username" ON "user" ("username");
CREATE INDEX "user_index_created_at" ON "user" ("created_at");
CREATE TRIGGER update_user_timestamp
  BEFORE UPDATE ON "user"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- category
CREATE TABLE "category" (
  "id" uuid NOT NULL,
  "label" varchar(100) NOT NULL UNIQUE,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  "search_vector" tsvector GENERATED ALWAYS AS (to_tsvector('english', label)) STORED,
  PRIMARY KEY ("id")
);
CREATE INDEX "category_index_label" ON "category" ("label");
CREATE INDEX "category_index_created_at" ON "category" ("created_at");
CREATE INDEX "category_search_vector_idx" ON "category" USING GIN ("search_vector");
CREATE TRIGGER update_category_timestamp
  BEFORE UPDATE ON "category"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- ingredient
CREATE TABLE "ingredient" (
  "id" uuid NOT NULL,
  "name" varchar(100) NOT NULL UNIQUE,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  "search_vector" tsvector GENERATED ALWAYS AS (to_tsvector('english', name)) STORED,
  PRIMARY KEY ("id")
);
CREATE INDEX "ingredient_index_name" ON "ingredient" ("name");
CREATE INDEX "ingredient_index_created_at" ON "ingredient" ("created_at");
CREATE INDEX "ingredient_search_vector_idx" ON "ingredient" USING GIN ("search_vector");
CREATE TRIGGER update_ingredient_timestamp
  BEFORE UPDATE ON "ingredient"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- tag
CREATE TABLE "tag" (
  "id" uuid NOT NULL,
  "name" varchar(50) NOT NULL UNIQUE,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "tag_index_name" ON "tag" ("name");
CREATE INDEX "tag_index_created_at" ON "tag" ("created_at");
CREATE TRIGGER update_tag_timestamp
  BEFORE UPDATE ON "tag"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- recipe
CREATE TABLE "recipe" (
  "id" uuid NOT NULL,
  "title" varchar(255) NOT NULL,
  "category_id" uuid NOT NULL,
  "creator_id" uuid NOT NULL,
  "preparation_time" bigint NOT NULL,
  "thumbnail_id" uuid,
  "like_count" bigint DEFAULT 0,
  "rating_count" bigint DEFAULT 0,
  "average_rating" decimal(3,2) DEFAULT 0.00,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  "deleted_at" timestamp with time zone,
  "search_vector" tsvector GENERATED ALWAYS AS (to_tsvector('english', title)) STORED,
  PRIMARY KEY ("id")
);
CREATE INDEX "recipe_index_title" ON "recipe" ("title");
CREATE INDEX "recipe_index_category_id" ON "recipe" ("category_id");
CREATE INDEX "recipe_index_creator_id" ON "recipe" ("creator_id");
CREATE INDEX "recipe_index_created_at" ON "recipe" ("created_at");
CREATE INDEX "recipe_index_preparation_time" ON "recipe" ("preparation_time");
CREATE INDEX "recipe_search_vector_idx" ON "recipe" USING GIN ("search_vector");
CREATE TRIGGER update_recipe_timestamp
  BEFORE UPDATE ON "recipe"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- recipe_tag
CREATE TABLE "recipe_tag" (
  "recipe_id" uuid NOT NULL,
  "tag_id" uuid NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("recipe_id", "tag_id")
);
CREATE INDEX "recipe_tag_index_recipe_id" ON "recipe_tag" ("recipe_id");
CREATE INDEX "recipe_tag_index_tag_id" ON "recipe_tag" ("tag_id");

-- recipe_step
CREATE TABLE "recipe_step" (
  "id" uuid NOT NULL,
  "recipe_id" uuid NOT NULL,
  "index" integer NOT NULL,
  "description" text NOT NULL,
  "picture_id" uuid,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "recipe_step_index_recipe_step" ON "recipe_step" ("recipe_id", "index");
CREATE TRIGGER update_recipe_step_timestamp
  BEFORE UPDATE ON "recipe_step"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- recipe_picture
CREATE TABLE "recipe_picture" (
  "id" uuid NOT NULL,
  "recipe_id" uuid NOT NULL,
  "path" varchar(255) NOT NULL UNIQUE,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("id")
);
CREATE INDEX "recipe_picture_index_path" ON "recipe_picture" ("path");
CREATE INDEX "recipe_picture_index_created_at" ON "recipe_picture" ("created_at");
CREATE INDEX "recipe_picture_index_recipe_id" ON "recipe_picture" ("recipe_id");
CREATE TRIGGER update_recipe_picture_timestamp
  BEFORE UPDATE ON "recipe_picture"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- recipe_ingredient
CREATE TABLE "recipe_ingredient" (
  "recipe_id" uuid NOT NULL,
  "ingredient_id" uuid NOT NULL,
  "quantity" decimal NOT NULL,
  "unit" varchar(50) NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("recipe_id", "ingredient_id")
);
CREATE INDEX "recipe_ingredient_index_recipe_ingredient" ON "recipe_ingredient" ("recipe_id", "ingredient_id");
CREATE INDEX "recipe_ingredient_index_ingredient_created_at" ON "recipe_ingredient" ("ingredient_id", "created_at");

-- liked_recipe
CREATE TABLE "liked_recipe" (
  "recipe_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("recipe_id", "user_id")
);
CREATE INDEX "liked_recipe_index_user_id_created_at" ON "liked_recipe" ("user_id", "created_at");
CREATE INDEX "liked_recipe_index_recipe_id_created_at" ON "liked_recipe" ("recipe_id", "created_at");

-- rating
CREATE TABLE "rating" (
  "recipe_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "value" integer NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("recipe_id", "user_id"),
  CONSTRAINT "check_rating_value" CHECK ("value" >= 1 AND "value" <= 5)
);
CREATE INDEX "rating_index_user_id_created_at" ON "rating" ("user_id", "created_at");
CREATE INDEX "rating_index_recipe_id_created_at" ON "rating" ("recipe_id", "created_at");

-- bookmark
CREATE TABLE "bookmark" (
  "recipe_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  PRIMARY KEY ("recipe_id", "user_id")
);
CREATE INDEX "bookmark_index_user_id_created_at" ON "bookmark" ("user_id", "created_at");
CREATE INDEX "bookmark_index_recipe_id_created_at" ON "bookmark" ("recipe_id", "created_at");

-- comment
CREATE TABLE "comment" (
  "id" uuid NOT NULL,
  "recipe_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "content" text NOT NULL,
  "created_at" timestamp with time zone NOT NULL,
  "updated_at" timestamp with time zone NOT NULL,
  "deleted_at" timestamp with time zone,
  PRIMARY KEY ("id")
);
CREATE INDEX "comment_index_recipe_id" ON "comment" ("recipe_id");
CREATE INDEX "comment_index_user_id" ON "comment" ("user_id");
CREATE INDEX "comment_index_created_at" ON "comment" ("created_at");
CREATE TRIGGER update_comment_timestamp
  BEFORE UPDATE ON "comment"
  FOR EACH ROW
  EXECUTE FUNCTION update_timestamp();

-- Foreign Keys
ALTER TABLE "recipe" ADD CONSTRAINT "fk_recipe_category_id" FOREIGN KEY ("category_id") REFERENCES "category" ("id");
ALTER TABLE "recipe" ADD CONSTRAINT "fk_recipe_creator_id" FOREIGN KEY ("creator_id") REFERENCES "user" ("id");
ALTER TABLE "recipe" ADD CONSTRAINT "fk_recipe_thumbnail_id" FOREIGN KEY ("thumbnail_id") REFERENCES "recipe_picture" ("id");
ALTER TABLE "recipe_step" ADD CONSTRAINT "fk_recipe_step_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "recipe_step" ADD CONSTRAINT "fk_recipe_step_picture_id" FOREIGN KEY ("picture_id") REFERENCES "recipe_picture" ("id");
ALTER TABLE "recipe_picture" ADD CONSTRAINT "fk_recipe_picture_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "recipe_ingredient" ADD CONSTRAINT "fk_recipe_ingredient_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "recipe_ingredient" ADD CONSTRAINT "fk_recipe_ingredient_ingredient_id" FOREIGN KEY ("ingredient_id") REFERENCES "ingredient" ("id");
ALTER TABLE "recipe_tag" ADD CONSTRAINT "fk_recipe_tag_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "recipe_tag" ADD CONSTRAINT "fk_recipe_tag_tag_id" FOREIGN KEY ("tag_id") REFERENCES "tag" ("id");
ALTER TABLE "liked_recipe" ADD CONSTRAINT "fk_liked_recipe_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "liked_recipe" ADD CONSTRAINT "fk_liked_recipe_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "rating" ADD CONSTRAINT "fk_rating_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "rating" ADD CONSTRAINT "fk_rating_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "bookmark" ADD CONSTRAINT "fk_bookmark_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "bookmark" ADD CONSTRAINT "fk_bookmark_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id");
ALTER TABLE "comment" ADD CONSTRAINT "fk_comment_recipe_id" FOREIGN KEY ("recipe_id") REFERENCES "recipe" ("id");
ALTER TABLE "comment" ADD CONSTRAINT "fk_comment_user_id" FOREIGN KEY ("user_id") REFERENCES "user" ("id");

-- Trigger for like_count
CREATE OR REPLACE FUNCTION update_like_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE "recipe"
        SET like_count = like_count + 1
        WHERE id = NEW.recipe_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE "recipe"
        SET like_count = GREATEST(like_count - 1, 0)
        WHERE id = OLD.recipe_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_like_count
AFTER INSERT OR DELETE ON "liked_recipe"
FOR EACH ROW EXECUTE FUNCTION update_like_count();

-- Trigger function to increment/decrement rating_count and update average_rating
CREATE OR REPLACE FUNCTION update_rating_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE "recipe"
        SET rating_count = rating_count + 1,
            average_rating = (
                CASE
                    WHEN rating_count = 0 THEN NEW.value
                    ELSE ((average_rating * rating_count + NEW.value) / (rating_count + 1))::decimal(3,2)
                END
            )
        WHERE id = NEW.recipe_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE "recipe"
        SET rating_count = GREATEST(rating_count - 1, 0),
            average_rating = (
                CASE
                    WHEN rating_count <= 1 THEN 0.00
                    ELSE ((
                        (average_rating * rating_count - OLD.value) / (rating_count - 1)
                    ))::decimal(3,2)
                END
            )
        WHERE id = OLD.recipe_id;
    ELSIF TG_OP = 'UPDATE' THEN
        UPDATE "recipe"
        SET average_rating = (
            CASE
                WHEN rating_count = 0 THEN 0.00
                ELSE ((
                    (average_rating * rating_count - OLD.value + NEW.value) / rating_count
                ))::decimal(3,2)
            END
        )
        WHERE id = NEW.recipe_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_rating_stats
AFTER INSERT OR UPDATE OF value OR DELETE ON "rating"
FOR EACH ROW EXECUTE FUNCTION update_rating_stats();

CREATE FUNCTION update_search_vector() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := to_tsvector('english', NEW.title);
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvector_update_trigger
BEFORE INSERT OR UPDATE ON recipe
FOR EACH ROW
EXECUTE FUNCTION tsvector_update_trigger(search_vector, 'pg_catalog.english', title);
