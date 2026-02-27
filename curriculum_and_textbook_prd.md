## CosmicQ Curriculum & Textbook Platform – PRD

This PRD describes the **functional requirements and data model expectations**
for the CosmicQ curriculum + textbook system. It is written so another project
can implement a compatible schema and backend just from this document.

The goal: **given a learner’s context** (country, “board” or equivalent,
grade, medium, subject, and textbook), the system should:

- Know **which curriculum standard** applies.
- Use one or more **reference textbooks** that fully cover that standard.
- Optionally accept **user‑uploaded textbooks**.
- Serve AI‑generated learning experiences (explanations, questions, slides,
  etc.) grounded in the correct curriculum.

---

## 1. Scope & Goals

### 1.1 In scope

- Support **multiple countries** with very different education systems.
- Model:
  - Countries.
  - Curriculum authorities (called **boards** in this PRD).
  - Optional **states/regions** that map to boards (e.g. US).
  - Grade systems and grades.
  - Mediums of instruction (languages).
  - Subjects.
  - Textbooks (books) and their chapters.
  - AI‑generated content (slides, worksheets, questions, etc.).
- Allow **multiple textbooks per (grade, subject, medium)**:
  - At least one “reference book” per curriculum standard.
  - Optional user‑uploaded books aligned to the same standard.

### 1.2 Out of scope (for this PRD)

- Detailed auth, billing, analytics dashboards.
- Specific embedding or vector‑DB implementation.
- Recommendation algorithms internals.

---

## 2. Key Concepts (Domain Model)

### 2.1 Country

Represents a nation (India, US, etc.).

- Must support:
  - `country_code` (e.g. ISO code).
  - Human‑readable `title`.
  - Flags that drive UI:
    - Whether users select a **board** (e.g. India).
    - Whether users select a **state/region** (e.g. US).

### 2.2 Board (Curriculum Authority)

Represents a **curriculum framework / authority**, not an individual school.

- Examples:
  - India: CBSE, ICSE, state boards like GSEB.
  - US: “Texas TEKS”, “California Common Core variant”, “US K–12 Generic”.
- Responsibilities:
  - Ties to a **grade system** (grade method).
  - Owns a set of **subjects** via `(board, grade, medium)` combinations.
  - Acts as the main “curriculum anchor” for textbooks and generated content.

### 2.3 State / Region (Optional)

Some countries (e.g. US) are better modeled as:

- **Country → State** in the UI, where each state uses some **board**.
- Requirements:
  - A state belongs to exactly one country.
  - A state can be **mapped to a default board** (its curriculum).
  - Users in that state see subjects/books associated with that board.

### 2.4 Grade Method & Grade

**Grade Method**:

- Describes a **grade system** (e.g. “India 1–12”, “US K–12”).
- Each board uses one grade method.

**Grade**:

- A concrete level in a grade method (e.g. “Standard 8”, “Grade 7”).
- Required attributes:
  - Display name (`title`).
  - Ordered position (`display_order`).
  - Optional numeric equivalent (e.g. 8 for Std 8).
  - Optional academic stage (primary / middle / secondary) for UI grouping.

### 2.5 Medium (Language)

Represents the **medium of instruction**, typically a language.

- Examples: English, Gujarati, Hindi, Spanish.
- Requirements:
  - Belongs to a country.
  - Optionally associated to a specific board.

### 2.6 Subject

Represents a **subject within a specific curriculum context**:

- Identified by the tuple:
  - Country
  - Board (curriculum)
  - Medium
  - Grade
- Examples:
  - India, CBSE, English, Std 8, Mathematics.
  - US, Texas TEKS, English, Grade 7, Math.
- Every textbook is attached to exactly one subject.

### 2.7 Book (Textbook)

Represents a **textbook or similar resource file**.

- Each book:
  - Belongs to a **subject**.
  - Also references **grade** and **medium** explicitly.
  - Has metadata:
    - `book_type` (e.g. reference, user_upload, supplemental).
    - Title, author, publisher, edition, year.
    - File paths (original, processed).
    - Curriculum version and effective dates.
    - Visibility / status flags.
    - Uploader information (for user books).
- Multiple books may exist per `(subject, grade, medium)`:
  - At least **one reference book** that fully covers the standard.
  - Optional **user‑uploaded** books for personalization.

### 2.8 Chapter

Represents a **chapter or logical section** of a book.

- Key attributes:
  - `book_id`.
  - Chapter number and title.
  - Page range (optional).
  - Pedagogical metadata:
    - Learning objectives.
    - Key concepts.
    - Difficulty.
    - Concept tags.
  - Embedding / AI metadata (if needed).

### 2.9 Generated Content

Represents **AI‑generated artifacts** linked back to the curriculum:

- Examples: slides, worksheets, quizzes, lesson plans.
- Must reference:
  - `chapter_id`.
  - `subject_id`.
  - `grade_id`.
  - `medium_id`.
  - `board_id`.
  - Optional `state_id` (if states are modeled separately).
- Also includes:
  - Usage statistics.
  - Ratings / quality metrics.
  - Sharing scope (private, school, global, etc.).

### 2.10 User Context

Tracks a user’s **current curriculum selection**:

- `current_country_id`
- `current_state_id` (optional)
- `current_board_id`
- `current_medium_id`
- `current_grade_id`
- `current_subject_id`

This allows the backend to resolve “where the user is” in the curriculum at any
time and fetch appropriate subjects, books, and content.

---

## 3. Core User Flows (Functional Requirements)

### 3.1 Country & curriculum selection

**User story:** As a teacher or student, I want to choose my country and
curriculum context so that all content is aligned to my syllabus.

Requirements:

1. User selects **country**.
2. Depending on the country configuration:
   - If `have_board = TRUE`: user sees a **board selector**.
   - If `has_states = TRUE`: user sees a **state selector**, and the system
     internally resolves that to a board.
3. User then selects:
   - Medium (language).
   - Grade.
   - Subject.
4. System updates `user_context` with the selected IDs.
5. All subsequent textbook and content queries use those context values.

### 3.2 Admin: configuring countries, boards, and states

**User story:** As an admin, I want to configure countries, boards, and states
so the UI can offer the correct curriculum choices.

Requirements:

- Admin can:
  - Create/edit countries with flags (`have_board`, `has_states`).
  - Create/edit boards linked to countries and grade methods.
  - (If using states) Create/edit states, each mapped to a **default board**.
- For US‑style systems:
  - Admin sets, for each state, which board (curriculum) it uses.
  - The user only chooses the state; the board is inferred.

### 3.3 Admin: seeding grade systems, grades, mediums, and subjects

**User story:** As an admin, I want to define grades, mediums, and subjects for
each board so that teachers and students can select them.

Requirements:

- Admin can:
  - Create grade methods and grades.
  - Assign a grade method to each board.
  - Create mediums per country (and optionally per board).
  - For each relevant combination of `(country, board, grade, medium)`:
    - Create subject rows (Math, Science, etc.).

### 3.4 Admin/Teacher: uploading reference textbooks

**User story:** As an admin or designated teacher, I want to upload reference
textbooks that fully cover a curriculum standard, so AI can use them internally
to teach.

Requirements:

- For a given `(country, board, grade, subject, medium)`:
  - System must support uploading **at least one “reference” book** that:
    - Is marked as covering the full curriculum for that standard.
- The system must:
  - Store the original file path and processed version.
  - Allow multiple revisions (via status or versioning).

### 3.5 Teacher/Student: uploading personal textbooks

**User story:** As a teacher or student, I want to upload my own textbook so
that AI can still teach me according to the official standard, but using my
book as context.

Requirements:

- Users can upload books and attach them to:
  - An existing `(country, board, grade, subject, medium)` context.
- These books must:
  - Be stored as separate `book` rows with a different `book_type` (e.g.
    `user_upload`).
  - Not change the definition of the official standard.
- AI behavior:
  - When generating explanations, it should:
    - Stay aligned to the **standard** implied by board + grade + subject.
    - Optionally use user book content for examples, exercises, etc.

### 3.6 Generated content production

**User story:** As a teacher, I want AI to generate content (slides, questions,
etc.) that is clearly linked to the correct curriculum context.

Requirements:

- Every generated content item must store:
  - Curriculum identifiers: `board_id`, `grade_id`, `subject_id`, `medium_id`.
  - Source chapter (`chapter_id`).
  - Optional `state_id` for localization/analytics.
- The system must support:
  - Filtering generated content by curriculum context.
  - Tracking usage and user feedback (accept/reject, ratings).

---

## 4. Data Model Expectations (for any implementation)

Any implementation based on this PRD should satisfy:

1. **Normalized core entities**:
   - Separate tables for: countries, boards, grade_methods, grades, mediums,
     subjects, books, chapters, users, user_context, generated_content.
   - Optional: states (if not using the simpler “state = board” model).

2. **Stable foreign‑key graph**:
   - Each entity’s relationships should enforce:
     - `board.country_id` → countries.
     - `grade.grade_method_id` → grade_methods.
     - `board.grade_method_id` → grade_methods.
     - `subject` references `country`, `board`, `grade`, `medium`.
     - `book` references `subject`, `grade`, `medium`.
     - `chapter` references `book`.
     - `generated_content` references at least:
       - `chapter`, `subject`, `grade`, `medium`, `board`.
     - `user_context` references the active curriculum selection.

3. **Efficient querying**:
   - Indexes that make these core queries fast:
     - Get subjects for `(country, board, medium, grade)`.
     - Get books for `(subject, medium, grade)`.
     - Get generated content for `(board, grade, subject, medium)` (+ optional
       filters like `state` or `chapter`).

4. **Curriculum‑first design**:
   - Curriculum standard is always determined by:
     - Country
     - Board (directly, or via state → default_board)
     - Grade
     - Subject
     - Medium
   - Books and generated content are always tied back to that context.

---

## 5. Variations & Implementation Notes

### 5.1 “State = board” simplification

For some projects it may be acceptable to **represent states directly as
boards**:

- US: each state is a `board` row.
- UI labels them as “State” for that country.
- This removes the need for a separate `states` table and `default_board_id`.

This PRD supports both:

- **Strict model**: separate `states` that map to boards.
- **Simplified model**: only `boards`, where some boards represent states.

The rest of the data model (subjects, books, chapters, generated content)
remains the same.

### 5.2 Multi‑board states (future)

If a future system needs a single state to support **multiple curricula** (e.g.
public vs private tracks), it can:

- Either keep multiple `board` rows and allow the user to pick board after
  state; or
- Introduce a **many‑to‑many relation** between states and boards.

This is considered a future extension and **not required** for an MVP.

---

## 6. Success Criteria

An implementation is considered compliant with this PRD if:

1. For any supported country, a user can:
   - Select their curriculum context using:
     - Country (+ board) or
     - Country + state (mapped to a board),
     - then medium, grade, subject.
2. Admins can:
   - Seed boards, grades, mediums, subjects, and books per curriculum.
3. The system can:
   - Store multiple textbooks per `(board, grade, subject, medium)`.
   - Distinguish reference books vs user uploads.
   - Attach generated content to a precise curriculum context.
4. All curriculum relationships are explicit in the data model, so:
   - A separate project can recreate a compatible schema by following:
     - Entity list and relationships in Sections 2 & 4.

This PRD is intentionally **implementation‑agnostic**: any SQL or NoSQL schema
that cleanly satisfies these entities, relationships, and flows is acceptable.

