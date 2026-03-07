-- =============================================================================
-- SEED DATA — India (Gujarat), United States, UAE
-- Run AFTER run_all_migrations.sql
--
-- HOW TO RUN in DBeaver:
--   1. Open this file
--   2. Make sure correct DB is selected in the connection dropdown
--   3. Press Ctrl+Alt+X  (Execute Script — NOT the plain Run button)
--
-- Every execution clears all data first, then re-inserts fresh data.
-- =============================================================================

-- ─────────────────────────────────────────────────────────────────────────────
-- CLEAR ALL DATA (reverse FK order) + reset sequences
-- ─────────────────────────────────────────────────────────────────────────────
TRUNCATE TABLE
    generated_content,
    chapters,
    books,
    user_default,
    users,
    subjects,
    mediums,
    grades,
    boards,
    languages,
    grade_methods,
    countries
RESTART IDENTITY CASCADE;

-- ─────────────────────────────────────────────────────────────────────────────
-- 1. COUNTRIES
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO countries (id, country_code, title, phone_code, signup_methods, have_board, is_visible)
VALUES
  (1, 'IN', 'India',         '+91',  ARRAY['email','phone'], TRUE, TRUE),
  (2, 'US', 'United States', '+1',   ARRAY['email'],         TRUE, TRUE),
  (3, 'AE', 'UAE',           '+971', ARRAY['email','phone'], TRUE, TRUE);

SELECT setval('countries_id_seq', (SELECT MAX(id) FROM countries));

-- ─────────────────────────────────────────────────────────────────────────────
-- 2. GRADE METHODS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO grade_methods (id, title, description, is_visible)
VALUES
  (1, 'India K-12', 'Indian school grading system (Class 1-12)',           TRUE),
  (2, 'US K-12',    'United States K-12 grading system (Kindergarten-12)', TRUE),
  (3, 'UAE K-12',   'UAE Ministry of Education K-12 grading system',       TRUE);

SELECT setval('grade_methods_id_seq', (SELECT MAX(id) FROM grade_methods));

-- ─────────────────────────────────────────────────────────────────────────────
-- 3. LANGUAGES
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO languages (id, code, name, is_visible)
VALUES
  (1, 'gu', 'Gujarati', TRUE),
  (2, 'hi', 'Hindi',    TRUE),
  (3, 'en', 'English',  TRUE),
  (4, 'ar', 'Arabic',   TRUE);

SELECT setval('languages_id_seq', (SELECT MAX(id) FROM languages));

-- ─────────────────────────────────────────────────────────────────────────────
-- 4. BOARDS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO boards (id, country_id, title, grade_method_id, is_visible)
VALUES
  (1, 1, 'GSEB',        1, TRUE),
  (2, 1, 'CBSE',        1, TRUE),
  (3, 2, 'Common Core', 2, TRUE),
  (4, 3, 'KHDA',        3, TRUE);

SELECT setval('boards_id_seq', (SELECT MAX(id) FROM boards));

-- ─────────────────────────────────────────────────────────────────────────────
-- 5. GRADES
-- India K-12  (grade_method_id=1) -> ids  1-12
-- US K-12     (grade_method_id=2) -> ids 13-25  (13=Kindergarten)
-- UAE K-12    (grade_method_id=3) -> ids 26-37
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO grades (id, grade_method_id, title, age_range_start, age_range_end, is_visible)
VALUES
  -- India
  ( 1, 1, 'Class 1',  6,  7,  TRUE), ( 2, 1, 'Class 2',  7,  8,  TRUE),
  ( 3, 1, 'Class 3',  8,  9,  TRUE), ( 4, 1, 'Class 4',  9,  10, TRUE),
  ( 5, 1, 'Class 5',  10, 11, TRUE), ( 6, 1, 'Class 6',  11, 12, TRUE),
  ( 7, 1, 'Class 7',  12, 13, TRUE), ( 8, 1, 'Class 8',  13, 14, TRUE),
  ( 9, 1, 'Class 9',  14, 15, TRUE), (10, 1, 'Class 10', 15, 16, TRUE),
  (11, 1, 'Class 11', 16, 17, TRUE), (12, 1, 'Class 12', 17, 18, TRUE),
  -- US
  (13, 2, 'Kindergarten', 5,  6,  TRUE),
  (14, 2, 'Grade 1',  6,  7,  TRUE), (15, 2, 'Grade 2',  7,  8,  TRUE),
  (16, 2, 'Grade 3',  8,  9,  TRUE), (17, 2, 'Grade 4',  9,  10, TRUE),
  (18, 2, 'Grade 5',  10, 11, TRUE), (19, 2, 'Grade 6',  11, 12, TRUE),
  (20, 2, 'Grade 7',  12, 13, TRUE), (21, 2, 'Grade 8',  13, 14, TRUE),
  (22, 2, 'Grade 9',  14, 15, TRUE), (23, 2, 'Grade 10', 15, 16, TRUE),
  (24, 2, 'Grade 11', 16, 17, TRUE), (25, 2, 'Grade 12', 17, 18, TRUE),
  -- UAE
  (26, 3, 'Grade 1',  6,  7,  TRUE), (27, 3, 'Grade 2',  7,  8,  TRUE),
  (28, 3, 'Grade 3',  8,  9,  TRUE), (29, 3, 'Grade 4',  9,  10, TRUE),
  (30, 3, 'Grade 5',  10, 11, TRUE), (31, 3, 'Grade 6',  11, 12, TRUE),
  (32, 3, 'Grade 7',  12, 13, TRUE), (33, 3, 'Grade 8',  13, 14, TRUE),
  (34, 3, 'Grade 9',  14, 15, TRUE), (35, 3, 'Grade 10', 15, 16, TRUE),
  (36, 3, 'Grade 11', 16, 17, TRUE), (37, 3, 'Grade 12', 17, 18, TRUE);

SELECT setval('grades_id_seq', (SELECT MAX(id) FROM grades));

-- ─────────────────────────────────────────────────────────────────────────────
-- 6. MEDIUMS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO mediums (id, country_id, board_id, title, language_id, is_visible)
VALUES
  (1, 1, 1, 'Gujarati Medium', 1, TRUE),
  (2, 1, 1, 'English Medium',  3, TRUE),
  (3, 1, 2, 'Hindi Medium',    2, TRUE),
  (4, 1, 2, 'English Medium',  3, TRUE),
  (5, 2, 3, 'English Medium',  3, TRUE),
  (6, 3, 4, 'Arabic Medium',   4, TRUE),
  (7, 3, 4, 'English Medium',  3, TRUE);

SELECT setval('mediums_id_seq', (SELECT MAX(id) FROM mediums));

-- ─────────────────────────────────────────────────────────────────────────────
-- 7. SUBJECTS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO subjects (id, board_id, medium_id, grade_id, title, subject_code, is_visible)
VALUES
  -- GSEB Gujarati Medium - Class 1
  ( 1, 1, 1,  1, 'Gujarati',             'GUJ-G1',  TRUE),
  ( 2, 1, 1,  1, 'Mathematics',          'MAT-G1',  TRUE),
  ( 3, 1, 1,  1, 'Paryavaran (EVS)',      'EVS-G1',  TRUE),
  -- GSEB Gujarati Medium - Class 5
  ( 4, 1, 1,  5, 'Gujarati',             'GUJ-G5',  TRUE),
  ( 5, 1, 1,  5, 'Mathematics',          'MAT-G5',  TRUE),
  ( 6, 1, 1,  5, 'Science',              'SCI-G5',  TRUE),
  ( 7, 1, 1,  5, 'Social Science',       'SSC-G5',  TRUE),
  ( 8, 1, 1,  5, 'English',              'ENG-G5',  TRUE),
  -- GSEB Gujarati Medium - Class 8
  ( 9, 1, 1,  8, 'Gujarati',             'GUJ-G8',  TRUE),
  (10, 1, 1,  8, 'Mathematics',          'MAT-G8',  TRUE),
  (11, 1, 1,  8, 'Science',              'SCI-G8',  TRUE),
  (12, 1, 1,  8, 'Social Science',       'SSC-G8',  TRUE),
  (13, 1, 1,  8, 'English',              'ENG-G8',  TRUE),
  (14, 1, 1,  8, 'Sanskrit',             'SAN-G8',  TRUE),
  -- GSEB Gujarati Medium - Class 10
  (15, 1, 1, 10, 'Gujarati',             'GUJ-G10', TRUE),
  (16, 1, 1, 10, 'Mathematics',          'MAT-G10', TRUE),
  (17, 1, 1, 10, 'Science & Technology', 'SCI-G10', TRUE),
  (18, 1, 1, 10, 'Social Science',       'SSC-G10', TRUE),
  (19, 1, 1, 10, 'English',              'ENG-G10', TRUE),
  (20, 1, 1, 10, 'Sanskrit',             'SAN-G10', TRUE),
  -- GSEB English Medium - Class 10
  (21, 1, 2, 10, 'English',              'ENG-G10-EM', TRUE),
  (22, 1, 2, 10, 'Mathematics',          'MAT-G10-EM', TRUE),
  (23, 1, 2, 10, 'Science & Technology', 'SCI-G10-EM', TRUE),
  (24, 1, 2, 10, 'Social Science',       'SSC-G10-EM', TRUE),
  (25, 1, 2, 10, 'Gujarati',             'GUJ-G10-EM', TRUE),
  -- CBSE English Medium - Class 5
  (26, 2, 4,  5, 'English',              'CBSE-ENG-G5', TRUE),
  (27, 2, 4,  5, 'Mathematics',          'CBSE-MAT-G5', TRUE),
  (28, 2, 4,  5, 'Science',              'CBSE-SCI-G5', TRUE),
  (29, 2, 4,  5, 'Social Studies',       'CBSE-SST-G5', TRUE),
  (30, 2, 4,  5, 'Hindi',                'CBSE-HIN-G5', TRUE),
  -- CBSE English Medium - Class 10
  (31, 2, 4, 10, 'English',              'CBSE-ENG-G10', TRUE),
  (32, 2, 4, 10, 'Mathematics Standard', 'CBSE-MAT-G10', TRUE),
  (33, 2, 4, 10, 'Science',              'CBSE-SCI-G10', TRUE),
  (34, 2, 4, 10, 'Social Science',       'CBSE-SST-G10', TRUE),
  (35, 2, 4, 10, 'Hindi',                'CBSE-HIN-G10', TRUE),
  -- Common Core US - Grade 1
  (36, 3, 5, 14, 'English Language Arts', 'CC-ELA-G1', TRUE),
  (37, 3, 5, 14, 'Mathematics',           'CC-MAT-G1', TRUE),
  -- Common Core US - Grade 5
  (38, 3, 5, 18, 'English Language Arts', 'CC-ELA-G5', TRUE),
  (39, 3, 5, 18, 'Mathematics',           'CC-MAT-G5', TRUE),
  (40, 3, 5, 18, 'Science',               'CC-SCI-G5', TRUE),
  (41, 3, 5, 18, 'Social Studies',        'CC-SST-G5', TRUE),
  -- Common Core US - Grade 10
  (42, 3, 5, 23, 'English Language Arts', 'CC-ELA-G10', TRUE),
  (43, 3, 5, 23, 'Mathematics',           'CC-MAT-G10', TRUE),
  (44, 3, 5, 23, 'Biology',               'CC-BIO-G10', TRUE),
  (45, 3, 5, 23, 'US History',            'CC-HIS-G10', TRUE),
  (46, 3, 5, 23, 'Physical Education',    'CC-PE-G10',  TRUE),
  -- KHDA Arabic Medium - Grade 5
  (47, 4, 6, 30, 'Arabic',                'KHDA-AR-G5',  TRUE),
  (48, 4, 6, 30, 'Mathematics',           'KHDA-MAT-G5', TRUE),
  (49, 4, 6, 30, 'Science',               'KHDA-SCI-G5', TRUE),
  (50, 4, 6, 30, 'Islamic Studies',       'KHDA-ISL-G5', TRUE),
  (51, 4, 6, 30, 'Social Studies',        'KHDA-SST-G5', TRUE),
  -- KHDA English Medium - Grade 10
  (52, 4, 7, 35, 'English',               'KHDA-ENG-G10', TRUE),
  (53, 4, 7, 35, 'Mathematics',           'KHDA-MAT-G10', TRUE),
  (54, 4, 7, 35, 'Physics',               'KHDA-PHY-G10', TRUE),
  (55, 4, 7, 35, 'Chemistry',             'KHDA-CHE-G10', TRUE),
  (56, 4, 7, 35, 'Social Studies',        'KHDA-SST-G10', TRUE);

SELECT setval('subjects_id_seq', (SELECT MAX(id) FROM subjects));

-- ─────────────────────────────────────────────────────────────────────────────
-- 8. USERS
-- password = bcrypt of "Test@123456"
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO users (
    id, email, phone_number, password_hash,
    first_name, last_name, role,
    is_active, is_verified,
    preferable_subject, plateform_version
)
VALUES
  (1, 'admin@cosmic.edu', NULL,
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
   'Cosmic', 'Admin', 'admin', TRUE, TRUE,
   ARRAY['Mathematics','Science'], '1.0.0'),

  (2, 'ravi.patel@gseb.edu.in', '+919898001001',
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
   'Ravi', 'Patel', 'teacher', TRUE, TRUE,
   ARRAY['Mathematics','Science & Technology'], '1.0.0'),

  (3, 'sarah.johnson@school.us', NULL,
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
   'Sarah', 'Johnson', 'teacher', TRUE, TRUE,
   ARRAY['English Language Arts','Social Studies'], '1.0.0'),

  (4, 'ahmed.hassan@khda.edu.ae', '+971501234567',
   '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
   'Ahmed', 'Hassan', 'teacher', TRUE, TRUE,
   ARRAY['Mathematics','Physics'], '1.0.0');

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- ─────────────────────────────────────────────────────────────────────────────
-- 9. USER_DEFAULT
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO user_default (user_id, current_country_id, current_board_id, current_medium_id, current_grade_id)
VALUES
  (2, 1, 1, 1, 10),   -- Ravi:  India / GSEB / Gujarati / Class 10
  (3, 2, 3, 5, 23),   -- Sarah: US    / Common Core / English / Grade 10
  (4, 3, 4, 6, 30);   -- Ahmed: UAE   / KHDA / Arabic / Grade 5

-- ─────────────────────────────────────────────────────────────────────────────
-- 10. BOOKS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO books (id, subject_id, title, publisher, publication_year, created_by, is_public, is_active, file_path, total_pages, is_visible)
VALUES
  ( 1, 16, 'Ganit - Dhoran 10',                    'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class10/math.pdf',           280, TRUE),
  ( 2, 17, 'Vigyan ane Takniki - Dhoran 10',        'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class10/science.pdf',        320, TRUE),
  ( 3, 18, 'Samajik Vigyan - Dhoran 10',            'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class10/social_science.pdf', 240, TRUE),
  ( 4, 15, 'Gujarati Sahitya - Dhoran 10',          'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class10/gujarati.pdf',       180, TRUE),
  ( 5,  5, 'Ganit - Dhoran 5',                      'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class5/math.pdf',            120, TRUE),
  ( 6,  6, 'Paryavaran Vigyan - Dhoran 5',          'GCERT',           2024, 2, TRUE, TRUE, 'books/gseb/class5/evs.pdf',             140, TRUE),
  ( 7, 33, 'Science Textbook Class 10',             'NCERT',           2024, 1, TRUE, TRUE, 'books/cbse/class10/science.pdf',        295, TRUE),
  ( 8, 32, 'Mathematics Standard Class 10',         'NCERT',           2024, 1, TRUE, TRUE, 'books/cbse/class10/math.pdf',           300, TRUE),
  ( 9, 39, 'Mathematics Grade 5 - Go Math!',        'Houghton Mifflin',2023, 3, TRUE, TRUE, 'books/us/grade5/math.pdf',              420, TRUE),
  (10, 38, 'English Language Arts Grade 5',         'McGraw-Hill',     2023, 3, TRUE, TRUE, 'books/us/grade5/ela.pdf',               380, TRUE),
  (11, 43, 'Algebra II - Common Core',              'Pearson',         2023, 3, TRUE, TRUE, 'books/us/grade10/math.pdf',             460, TRUE),
  (12, 44, 'Biology - Exploring Life',              'Pearson',         2023, 3, TRUE, TRUE, 'books/us/grade10/bio.pdf',              510, TRUE),
  (13, 53, 'Mathematics Grade 10 - UAE Curriculum', 'ADEC',            2024, 4, TRUE, TRUE, 'books/uae/grade10/math.pdf',            350, TRUE),
  (14, 54, 'Physics Grade 10 - UAE Curriculum',     'ADEC',            2024, 4, TRUE, TRUE, 'books/uae/grade10/physics.pdf',         290, TRUE);

SELECT setval('books_id_seq', (SELECT MAX(id) FROM books));

-- ─────────────────────────────────────────────────────────────────────────────
-- 11. CHAPTERS
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO chapters (id, book_id, chapter_title, content_summary, concept_tags, is_visible)
VALUES
  -- Book 1: Ganit Dhoran 10 (GSEB Math Class 10)
  ( 1, 1, 'Vastav Sankhyao (Real Numbers)',
    'Euclid ni algorithm, Irrational numbers, HCF/LCM',
    ARRAY['real numbers','HCF','LCM','Euclid','irrational'], TRUE),
  ( 2, 1, 'Bahupadio (Polynomials)',
    'Zeroes of polynomials, relationship between zeroes and coefficients',
    ARRAY['polynomial','zeroes','coefficients','factorisation'], TRUE),
  ( 3, 1, 'Rekhiy Samikarano (Linear Equations)',
    'Pair of linear equations, graphical and algebraic methods',
    ARRAY['linear equations','substitution','elimination','graph'], TRUE),
  ( 4, 1, 'Tribhuj (Triangles)',
    'Similarity, Pythagoras theorem, area of similar triangles',
    ARRAY['triangle','similarity','Pythagoras','congruence'], TRUE),
  -- Book 2: Vigyan Dhoran 10 (GSEB Science Class 10)
  ( 5, 2, 'Rasayanik Prakriyao (Chemical Reactions)',
    'Types of reactions, oxidation-reduction, corrosion',
    ARRAY['chemical reaction','oxidation','reduction','acid','base'], TRUE),
  ( 6, 2, 'Tarango na Parkaar (Electricity)',
    'Ohm law, series and parallel circuits, electric power',
    ARRAY['Ohm law','current','resistance','circuit','power'], TRUE),
  ( 7, 2, 'Carbon ane tar Sanyojo (Carbon Compounds)',
    'Covalent bonds, homologous series, nomenclature',
    ARRAY['carbon','covalent','homologous','organic chemistry'], TRUE),
  -- Book 3: Samajik Vigyan Dhoran 10 (GSEB Social Science)
  ( 8, 3, 'Bharat: Sansadhan ane Vikas (Resources & Development)',
    'Types of resources, land use, soil types in India/Gujarat',
    ARRAY['resources','land use','soil','India','Gujarat','development'], TRUE),
  ( 9, 3, 'Lokshahi ane Vividhata (Democracy & Diversity)',
    'Social divisions, politics of social divisions in India',
    ARRAY['democracy','diversity','social division','Gujarat'], TRUE),
  (10, 3, 'Vastu ane Seva (Money & Credit)',
    'Money as a medium, formal and informal credit sources',
    ARRAY['money','credit','bank','SHG','formal credit'], TRUE),
  -- Book 7: CBSE Science Class 10
  (11, 7, 'Chemical Reactions and Equations',
    'Balancing equations, types of chemical reactions',
    ARRAY['chemical equation','balancing','combination','decomposition'], TRUE),
  (12, 7, 'Life Processes',
    'Nutrition, respiration, transportation, excretion',
    ARRAY['nutrition','respiration','photosynthesis','excretion'], TRUE),
  (13, 7, 'Heredity and Evolution',
    'Mendel laws, variation, natural selection',
    ARRAY['heredity','Mendel','gene','evolution','variation'], TRUE),
  -- Book 9: Go Math! Grade 5 US
  (14, 9, 'Place Value and Decimal Fractions',
    'Powers of 10, decimals to thousandths, comparing decimals',
    ARRAY['place value','decimal','thousandths','powers of 10'], TRUE),
  (15, 9, 'Multi-Digit Multiplication and Division',
    'Algorithms for multiplication and long division',
    ARRAY['multiplication','division','algorithm','multi-digit'], TRUE),
  (16, 9, 'Adding and Subtracting Fractions',
    'Unlike denominators, mixed numbers, word problems',
    ARRAY['fraction','denominator','mixed number','addition'], TRUE),
  -- Book 12: Biology US Grade 10
  (17, 12, 'Cell Biology',
    'Cell structure, organelles, cell membrane, mitosis',
    ARRAY['cell','organelle','mitosis','membrane','nucleus'], TRUE),
  (18, 12, 'Genetics',
    'DNA replication, transcription, translation, mutations',
    ARRAY['DNA','genetics','transcription','translation','mutation'], TRUE),
  (19, 12, 'Ecology',
    'Ecosystems, food chains, energy flow, biomes',
    ARRAY['ecosystem','food chain','energy flow','biome','population'], TRUE),
  -- Book 13: Mathematics UAE Grade 10
  (20, 13, 'Quadratic Equations',
    'Solving by factoring, quadratic formula, discriminant',
    ARRAY['quadratic','factoring','discriminant','roots'], TRUE),
  (21, 13, 'Trigonometry',
    'Sine, cosine, tangent, unit circle, identities',
    ARRAY['trigonometry','sine','cosine','unit circle','identity'], TRUE),
  (22, 13, 'Statistics and Probability',
    'Mean, median, mode, standard deviation, probability rules',
    ARRAY['statistics','probability','mean','standard deviation'], TRUE),
  -- Book 14: Physics UAE Grade 10
  (23, 14, 'Motion and Forces',
    'Newton laws, velocity, acceleration, free fall',
    ARRAY['Newton','force','velocity','acceleration','motion'], TRUE),
  (24, 14, 'Energy and Work',
    'Kinetic and potential energy, conservation of energy',
    ARRAY['energy','work','kinetic','potential','conservation'], TRUE),
  (25, 14, 'Waves and Sound',
    'Wave properties, sound waves, reflection, refraction',
    ARRAY['wave','sound','frequency','amplitude','refraction'], TRUE);

SELECT setval('chapters_id_seq', (SELECT MAX(id) FROM chapters));

-- ─────────────────────────────────────────────────────────────────────────────
-- 12. GENERATED_CONTENT
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO generated_content (
    id, content_type, chapter_id, generated_by_user_id,
    generation_prompt, generation_model,
    file_url, file_path, file_size_bytes, file_format,
    title, description, slide_count, question_count, concept_tags,
    medium_id, grade_id, subject_id, board_id,
    usage_count, download_count, view_count,
    recommendation_accept_count, recommendation_reject_count,
    quality_score, average_rating, rating_count,
    is_reusable, is_anonymous, is_public, share_scope, status
)
VALUES
  (1, 'presentation', 1, 2,
   'Create 15-slide presentation on Real Numbers for Class 10 GSEB Gujarati medium',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/1/real-numbers-slides.pdf',
   'generated/gseb/class10/math/real-numbers-slides.pdf',
   2048000, 'pdf',
   'Real Numbers - Class 10 GSEB',
   'Covers Euclid algorithm, HCF, LCM, irrational numbers in Gujarati context',
   15, NULL, ARRAY['real numbers','HCF','LCM','Euclid','irrational'],
   1, 10, 16, 1,
   45, 12, 200, 38, 4, 0.92, 4.50, 18,
   TRUE, FALSE, TRUE, 'global', 'active'),

  (2, 'quiz', 5, 2,
   'Generate 20 MCQ questions on Chemical Reactions for Class 10 GSEB Gujarati medium',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/2/chem-reactions-quiz.json',
   'generated/gseb/class10/science/chem-reactions-quiz.json',
   512000, 'json',
   'Chemical Reactions MCQ - Class 10 GSEB',
   '20 multiple-choice questions with answer key',
   NULL, 20, ARRAY['chemical reaction','oxidation','reduction','acid','base'],
   1, 10, 17, 1,
   88, 30, 350, 72, 6, 0.89, 4.30, 31,
   TRUE, FALSE, TRUE, 'global', 'active'),

  (3, 'presentation', 8, 2,
   'Create presentation on Resources and Development for Class 10 GSEB, Gujarat agriculture and land use focus',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/3/resources-dev-slides.pdf',
   'generated/gseb/class10/social/resources-dev-slides.pdf',
   1800000, 'pdf',
   'Resources & Development - Class 10 GSEB',
   'Gujarat context: Narmada, cotton belt, black soil, groundnut farming',
   12, NULL, ARRAY['resources','land use','soil','Gujarat','agriculture'],
   1, 10, 18, 1,
   33, 8, 120, 28, 2, 0.87, 4.40, 11,
   TRUE, FALSE, TRUE, 'global', 'active'),

  (4, 'quiz', 17, 3,
   'Generate 25 MCQ questions on Cell Biology for Grade 10 Common Core',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/4/cell-biology-quiz.json',
   'generated/us/grade10/bio/cell-biology-quiz.json',
   640000, 'json',
   'Cell Biology MCQ - Grade 10 US',
   '25 multiple-choice questions aligned with Common Core standards',
   NULL, 25, ARRAY['cell','organelle','mitosis','membrane'],
   5, 23, 44, 3,
   55, 20, 210, 48, 3, 0.91, 4.60, 22,
   TRUE, FALSE, TRUE, 'global', 'active'),

  (5, 'presentation', 20, 4,
   'Create 12-slide presentation on Quadratic Equations for Grade 10 KHDA curriculum, bilingual Arabic/English',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/5/quadratic-slides.pdf',
   'generated/uae/grade10/math/quadratic-slides.pdf',
   1600000, 'pdf',
   'Quadratic Equations - Grade 10 UAE',
   'KHDA-aligned, bilingual Arabic/English labels',
   12, NULL, ARRAY['quadratic','factoring','discriminant','roots'],
   7, 35, 53, 4,
   27, 7, 95, 22, 1, 0.88, 4.20, 9,
   TRUE, FALSE, TRUE, 'global', 'active'),

  (6, 'notes', 12, 1,
   'Generate concise revision notes on Life Processes for Class 10 CBSE',
   'claude-sonnet-4-6',
   'https://cdn.cosmic.edu/gc/6/life-processes-notes.pdf',
   'generated/cbse/class10/science/life-processes-notes.pdf',
   920000, 'pdf',
   'Life Processes - Revision Notes Class 10 CBSE',
   'Covers nutrition, respiration, transportation and excretion',
   NULL, NULL, ARRAY['nutrition','respiration','photosynthesis','excretion'],
   4, 10, 33, 2,
   60, 25, 310, 55, 3, 0.90, 4.50, 27,
   TRUE, FALSE, TRUE, 'global', 'active');

SELECT setval('generated_content_id_seq', (SELECT MAX(id) FROM generated_content));

-- =============================================================================
-- Done. All tables cleared and re-seeded.
-- =============================================================================
