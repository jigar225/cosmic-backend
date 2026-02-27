CREATE TABLE "countries"(
    "id" BIGINT NOT NULL,
    "country_code" VARCHAR(3) NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "icon_path" VARCHAR(255) NULL,
    "phone_code" VARCHAR(10) NULL,
    "signup_method" VARCHAR(50) NOT NULL DEFAULT 'email',
    "have_board" BOOLEAN NOT NULL DEFAULT TRUE,
    "has_states" BOOLEAN NOT NULL DEFAULT FALSE,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "countries" ADD PRIMARY KEY("id");
CREATE TABLE "grade_methods"(
    "id" BIGINT NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "description" TEXT NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "grade_methods" ADD PRIMARY KEY("id");
CREATE TABLE "boards"(
    "id" BIGINT NOT NULL,
    "country_id" BIGINT NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "grade_method_id" BIGINT NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "boards" ADD PRIMARY KEY("id");
CREATE INDEX "boards_country_id_index" ON
    "boards"("country_id");
CREATE TABLE "states"(
    "id" BIGINT NOT NULL,
    "country_id" BIGINT NOT NULL,
    "code" VARCHAR(10) NULL,
    "title" VARCHAR(100) NOT NULL,
    "default_board_id" BIGINT NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "states" ADD PRIMARY KEY("id");
CREATE INDEX "states_country_id_index" ON
    "states"("country_id");
CREATE INDEX "states_default_board_id_index" ON
    "states"("default_board_id");
CREATE TABLE "grades"(
    "id" BIGINT NOT NULL,
    "grade_method_id" BIGINT NOT NULL,
    "title" VARCHAR(50) NOT NULL,
    "display_order" INTEGER NOT NULL,
    "numeric_equivalent" INTEGER NULL,
    "age_range_start" INTEGER NULL,
    "age_range_end" INTEGER NULL,
    "academic_stage" VARCHAR(50) NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "grades" ADD PRIMARY KEY("id");
CREATE INDEX "grades_grade_method_id_index" ON
    "grades"("grade_method_id");
CREATE INDEX "grades_numeric_equivalent_index" ON
    "grades"("numeric_equivalent");
CREATE TABLE "mediums"(
    "id" BIGINT NOT NULL,
    "country_id" BIGINT NOT NULL,
    "board_id" BIGINT NULL,
    "title" VARCHAR(100) NOT NULL,
    "language_code" VARCHAR(10) NULL,
    "script_type" VARCHAR(50) NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "mediums" ADD PRIMARY KEY("id");
CREATE INDEX "mediums_country_id_index" ON
    "mediums"("country_id");
CREATE INDEX "mediums_board_id_index" ON
    "mediums"("board_id");
CREATE TABLE "subjects"(
    "id" BIGINT NOT NULL,
    "country_id" BIGINT NOT NULL,
    "board_id" BIGINT NOT NULL,
    "medium_id" BIGINT NOT NULL,
    "grade_id" BIGINT NOT NULL,
    "title" VARCHAR(100) NOT NULL,
    "subject_code" VARCHAR(50) NULL,
    "subject_type" VARCHAR(50) NOT NULL DEFAULT 'core',
    "sequence_order" INTEGER NULL,
    "description" TEXT NULL,
    "created_by" BIGINT NULL,
    "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "subjects_board_id_grade_id_index" ON
    "subjects"("board_id", "grade_id");
CREATE INDEX "subjects_country_id_board_id_medium_id_grade_id_index" ON
    "subjects"(
        "country_id",
        "board_id",
        "medium_id",
        "grade_id"
    );
ALTER TABLE
    "subjects" ADD PRIMARY KEY("id");
CREATE TABLE "users"(
    "id" BIGINT NOT NULL,
    "uuid" UUID NOT NULL DEFAULT UUID_GENERATE_V4(), "email" VARCHAR(255) NULL, "phone_number" VARCHAR(50) NULL, "password_hash" TEXT NOT NULL, "first_name" VARCHAR(100) NULL, "last_name" VARCHAR(100) NULL, "profile_photo" VARCHAR(500) NULL, "role" VARCHAR(50) NOT NULL DEFAULT 'teacher', "user_type" VARCHAR(50) NOT NULL DEFAULT 'individual', "is_active" BOOLEAN NOT NULL DEFAULT TRUE, "is_verified" BOOLEAN NULL, "email_verified_at" TIMESTAMP(0) WITH
        TIME zone NULL,
        "phone_verified_at" TIMESTAMP(0)
    WITH
        TIME zone NULL,
        "last_login_at" TIMESTAMP(0)
    WITH
        TIME zone NULL,
        "created_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "deleted_at" TIMESTAMP(0)
    WITH
        TIME zone NULL);
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "user_context"(
    "user_id" BIGINT NOT NULL,
    "current_country_id" BIGINT NULL,
    "current_state_id" BIGINT NULL,
    "current_board_id" BIGINT NULL,
    "current_medium_id" BIGINT NULL,
    "current_grade_id" BIGINT NULL,
    "current_subject_id" BIGINT NULL,
    "updated_at" TIMESTAMP(0) WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE
    "user_context" ADD PRIMARY KEY("user_id");
CREATE TABLE "books"(
    "id" BIGINT NOT NULL,
    "book_type" VARCHAR(50) NOT NULL,
    "subject_id" BIGINT NOT NULL,
    "medium_id" BIGINT NOT NULL,
    "grade_id" BIGINT NOT NULL,
    "title" VARCHAR(500) NOT NULL,
    "author" VARCHAR(255) NULL,
    "publisher" VARCHAR(255) NULL,
    "edition" VARCHAR(50) NULL,
    "publication_year" INTEGER NULL,
    "isbn" VARCHAR(20) NULL,
    "book_code" VARCHAR(100) NULL,
    "uploaded_by_user_id" BIGINT NULL,
    "is_public" BOOLEAN NULL,
    "curriculum_version" VARCHAR(50) NULL,
    "effective_start_date" DATE NULL,
    "effective_end_date" DATE NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'draft',
    "original_file_path" VARCHAR(1000) NOT NULL,
    "processed_file_path" VARCHAR(1000) NULL,
    "cover_image_url" VARCHAR(500) NULL,
    "total_pages" INTEGER NULL,
    "has_toc" BOOLEAN NULL,
    "toc_extraction_method" VARCHAR(50) NULL,
    "processing_status" VARCHAR(50) NULL,
    "processing_notes" TEXT NULL,
    "processing_started_at" TIMESTAMP(0) WITH
        TIME zone NULL,
        "processing_completed_at" TIMESTAMP(0)
    WITH
        TIME zone NULL,
        "description" TEXT NULL,
        "tags" TEXT[] NULL,
        "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
        "view_count" INTEGER NULL,
        "download_count" INTEGER NULL,
        "created_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "books_subject_id_medium_id_grade_id_index" ON
    "books"(
        "subject_id",
        "medium_id",
        "grade_id"
    );
ALTER TABLE
    "books" ADD PRIMARY KEY("id");
CREATE INDEX "books_book_type_index" ON
    "books"("book_type");
CREATE INDEX "books_status_index" ON
    "books"("status");
CREATE TABLE "chapters"(
    "id" BIGINT NOT NULL,
    "book_id" BIGINT NOT NULL,
    "chapter_number" INTEGER NOT NULL,
    "chapter_title" VARCHAR(500) NOT NULL,
    "chapter_code" VARCHAR(100) NULL,
    "page_start" INTEGER NULL,
    "page_end" INTEGER NULL,
    "section_type" VARCHAR(50) NULL,
    "content_summary" TEXT NULL,
    "word_count" INTEGER NULL,
    "learning_objectives" TEXT[] NULL,
    "key_concepts" TEXT[] NULL,
    "estimated_teaching_hours" DECIMAL(4, 2) NULL,
    "difficulty_level" VARCHAR(50) NULL,
    "prerequisites" TEXT[] NULL,
    "concept_tags" TEXT[] NULL,
    "concept_extraction_confidence" DECIMAL(3, 2) NULL,
    "concept_extraction_method" VARCHAR(50) NULL,
    "concept_verified" BOOLEAN NULL,
    "concept_verified_by" BIGINT NULL,
    "concept_verified_at" TIMESTAMP(0) WITH
        TIME zone NULL,
        "embedding_id" VARCHAR(100) NULL,
        "embedding_generated_at" TIMESTAMP(0)
    WITH
        TIME zone NULL,
        "embedding_model" VARCHAR(50) NULL,
        "embedding_dimensions" INTEGER NULL,
        "status" VARCHAR(50) NOT NULL DEFAULT 'draft',
        "version" INTEGER NOT NULL DEFAULT 1,
        "parent_chapter_id" BIGINT NULL,
        "previous_chapter_id" BIGINT NULL,
        "next_chapter_id" BIGINT NULL,
        "is_visible" BOOLEAN NOT NULL DEFAULT TRUE,
        "view_count" INTEGER NULL,
        "created_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "chapters_book_id_chapter_number_index" ON
    "chapters"("book_id", "chapter_number");
ALTER TABLE
    "chapters" ADD PRIMARY KEY("id");
CREATE INDEX "chapters_book_id_index" ON
    "chapters"("book_id");
CREATE INDEX "chapters_concept_tags_index" ON
    "chapters"("concept_tags");
CREATE INDEX "chapters_status_index" ON
    "chapters"("status");
CREATE TABLE "generated_content"(
    "id" BIGINT NOT NULL,
    "content_type" VARCHAR(50) NOT NULL,
    "chapter_id" BIGINT NOT NULL,
    "generated_by_user_id" BIGINT NOT NULL,
    "generation_prompt" TEXT NULL,
    "generation_params" jsonb NULL,
    "generation_model" VARCHAR(50) NULL,
    "generation_duration_seconds" INTEGER NULL,
    "file_url" VARCHAR(1000) NOT NULL,
    "file_path" VARCHAR(1000) NULL,
    "file_size_bytes" BIGINT NULL,
    "file_format" VARCHAR(20) NULL,
    "thumbnail_url" VARCHAR(500) NULL,
    "title" VARCHAR(255) NULL,
    "description" TEXT NULL,
    "slide_count" INTEGER NULL,
    "page_count" INTEGER NULL,
    "question_count" INTEGER NULL,
    "embedding_id" VARCHAR(100) NULL,
    "embedding_generated_at" TIMESTAMP(0) WITH
        TIME zone NULL,
        "embedding_model" VARCHAR(50) NULL,
        "embedding_dimensions" INTEGER NULL,
        "concept_tags" TEXT[] NULL,
        "medium_id" BIGINT NOT NULL,
        "grade_id" BIGINT NOT NULL,
        "subject_id" BIGINT NOT NULL,
    "board_id" BIGINT NOT NULL,
    "state_id" BIGINT NULL,
    "usage_count" INTEGER NOT NULL,
        "download_count" INTEGER NOT NULL,
        "view_count" INTEGER NOT NULL,
        "recommendation_accept_count" INTEGER NOT NULL,
        "recommendation_reject_count" INTEGER NOT NULL,
        "quality_score" DECIMAL(3, 2) NULL,
        "average_rating" DECIMAL(3, 2) NULL,
        "rating_count" INTEGER NOT NULL,
        "is_reusable" BOOLEAN NOT NULL DEFAULT TRUE,
        "is_anonymous" BOOLEAN NOT NULL DEFAULT TRUE,
        "is_public" BOOLEAN NOT NULL DEFAULT TRUE,
        "share_scope" VARCHAR(50) NOT NULL DEFAULT 'global',
        "status" VARCHAR(50) NOT NULL DEFAULT 'active',
        "created_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "updated_at" TIMESTAMP(0)
    WITH
        TIME zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
        "last_used_at" TIMESTAMP(0)
    WITH
        TIME zone NULL
);
CREATE INDEX "generated_content_medium_id_grade_id_subject_id_index" ON
    "generated_content"(
        "medium_id",
        "grade_id",
        "subject_id"
    );
CREATE INDEX "generated_content_quality_score_usage_count_index" ON
    "generated_content"("quality_score", "usage_count");
ALTER TABLE
    "generated_content" ADD PRIMARY KEY("id");
CREATE INDEX "generated_content_content_type_index" ON
    "generated_content"("content_type");
CREATE INDEX "generated_content_chapter_id_index" ON
    "generated_content"("chapter_id");
CREATE INDEX "generated_content_generated_by_user_id_index" ON
    "generated_content"("generated_by_user_id");
CREATE INDEX "generated_content_concept_tags_index" ON
    "generated_content"("concept_tags");
CREATE INDEX "generated_content_state_id_index" ON
    "generated_content"("state_id");
ALTER TABLE
    "subjects" ADD CONSTRAINT "subjects_board_id_foreign" FOREIGN KEY("board_id") REFERENCES "boards"("id");
ALTER TABLE
    "chapters" ADD CONSTRAINT "chapters_parent_chapter_id_foreign" FOREIGN KEY("parent_chapter_id") REFERENCES "chapters"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_board_id_foreign" FOREIGN KEY("current_board_id") REFERENCES "boards"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "books" ADD CONSTRAINT "books_subject_id_foreign" FOREIGN KEY("subject_id") REFERENCES "subjects"("id");
ALTER TABLE
    "mediums" ADD CONSTRAINT "mediums_country_id_foreign" FOREIGN KEY("country_id") REFERENCES "countries"("id");
ALTER TABLE
    "subjects" ADD CONSTRAINT "subjects_grade_id_foreign" FOREIGN KEY("grade_id") REFERENCES "grades"("id");
ALTER TABLE
    "boards" ADD CONSTRAINT "boards_country_id_foreign" FOREIGN KEY("country_id") REFERENCES "countries"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_subject_id_foreign" FOREIGN KEY("current_subject_id") REFERENCES "subjects"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_board_id_foreign" FOREIGN KEY("board_id") REFERENCES "boards"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_chapter_id_foreign" FOREIGN KEY("chapter_id") REFERENCES "chapters"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_medium_id_foreign" FOREIGN KEY("current_medium_id") REFERENCES "mediums"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_generated_by_user_id_foreign" FOREIGN KEY("generated_by_user_id") REFERENCES "users"("id");
ALTER TABLE
    "chapters" ADD CONSTRAINT "chapters_previous_chapter_id_foreign" FOREIGN KEY("previous_chapter_id") REFERENCES "chapters"("id");
ALTER TABLE
    "subjects" ADD CONSTRAINT "subjects_country_id_foreign" FOREIGN KEY("country_id") REFERENCES "countries"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_subject_id_foreign" FOREIGN KEY("subject_id") REFERENCES "subjects"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_country_id_foreign" FOREIGN KEY("current_country_id") REFERENCES "countries"("id");
ALTER TABLE
    "books" ADD CONSTRAINT "books_uploaded_by_user_id_foreign" FOREIGN KEY("uploaded_by_user_id") REFERENCES "users"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_grade_id_foreign" FOREIGN KEY("current_grade_id") REFERENCES "grades"("id");
ALTER TABLE
    "grades" ADD CONSTRAINT "grades_grade_method_id_foreign" FOREIGN KEY("grade_method_id") REFERENCES "grade_methods"("id");
ALTER TABLE
    "subjects" ADD CONSTRAINT "subjects_medium_id_foreign" FOREIGN KEY("medium_id") REFERENCES "mediums"("id");
ALTER TABLE
    "books" ADD CONSTRAINT "books_medium_id_foreign" FOREIGN KEY("medium_id") REFERENCES "mediums"("id");
ALTER TABLE
    "books" ADD CONSTRAINT "books_grade_id_foreign" FOREIGN KEY("grade_id") REFERENCES "grades"("id");
ALTER TABLE
    "mediums" ADD CONSTRAINT "mediums_board_id_foreign" FOREIGN KEY("board_id") REFERENCES "boards"("id");
ALTER TABLE
    "boards" ADD CONSTRAINT "boards_grade_method_id_foreign" FOREIGN KEY("grade_method_id") REFERENCES "grade_methods"("id");
ALTER TABLE
    "chapters" ADD CONSTRAINT "chapters_concept_verified_by_foreign" FOREIGN KEY("concept_verified_by") REFERENCES "users"("id");
ALTER TABLE
    "chapters" ADD CONSTRAINT "chapters_next_chapter_id_foreign" FOREIGN KEY("next_chapter_id") REFERENCES "chapters"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_medium_id_foreign" FOREIGN KEY("medium_id") REFERENCES "mediums"("id");
ALTER TABLE
    "chapters" ADD CONSTRAINT "chapters_book_id_foreign" FOREIGN KEY("book_id") REFERENCES "books"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_grade_id_foreign" FOREIGN KEY("grade_id") REFERENCES "grades"("id");
ALTER TABLE
    "states" ADD CONSTRAINT "states_country_id_foreign" FOREIGN KEY("country_id") REFERENCES "countries"("id");
ALTER TABLE
    "states" ADD CONSTRAINT "states_default_board_id_foreign" FOREIGN KEY("default_board_id") REFERENCES "boards"("id");
ALTER TABLE
    "user_context" ADD CONSTRAINT "user_context_current_state_id_foreign" FOREIGN KEY("current_state_id") REFERENCES "states"("id");
ALTER TABLE
    "generated_content" ADD CONSTRAINT "generated_content_state_id_foreign" FOREIGN KEY("state_id") REFERENCES "states"("id");