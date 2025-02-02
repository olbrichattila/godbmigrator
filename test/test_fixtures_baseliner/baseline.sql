CREATE TABLE migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME
	);
CREATE TABLE migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT
	);
CREATE TABLE users (
    name VARCHAR(255)
);
CREATE TABLE users2 (
    name VARCHAR(255)
);
CREATE TABLE users3 (
    name VARCHAR(255)
);
CREATE TABLE users4 (
    name VARCHAR(255)
);
CREATE TABLE users5 (
    name VARCHAR(255)
);
CREATE TABLE users6 (
    name VARCHAR(255)
);
CREATE TABLE olb_migrations (
		file_name VARCHAR(255),
		created_at DATETIME,
		deleted_at DATETIME);
CREATE TABLE olb_migration_reports (
		file_name VARCHAR(255),
		result_status VARCHAR(12),
		created_at DATETIME,
		message TEXT);
CREATE TABLE books (
    book_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER,
    genre_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES authors(author_id),
    FOREIGN KEY (genre_id) REFERENCES genres(genre_id)
);
CREATE TABLE members (
    member_id INTEGER PRIMARY KEY,
    member_name TEXT NOT NULL
);
CREATE TABLE authors (
    author_id INTEGER PRIMARY KEY,
    author_name TEXT NOT NULL
);
CREATE TABLE genres (
    genre_id INTEGER PRIMARY KEY,
    genre_name TEXT NOT NULL
);
CREATE TABLE loans (
    loan_id INTEGER PRIMARY KEY,
    book_id INTEGER,
    member_id INTEGER,
    loan_date DATE,
    return_date DATE,
    FOREIGN KEY (book_id) REFERENCES books(book_id),
    FOREIGN KEY (member_id) REFERENCES members(member_id)
);
CREATE INDEX idx_book_title ON books (title);
CREATE INDEX idx_book_author_title ON books (author_id, title);
DELIMITER ;
CREATE VIEW book_authors AS
SELECT b.title, a.author_name
FROM books AS b
JOIN authors AS a ON b.author_id = a.author_id
DELIMITER ;;
DELIMITER ;
CREATE TRIGGER check_loan_dates
BEFORE UPDATE OF return_date ON loans
FOR EACH ROW
WHEN NEW.return_date < OLD.loan_date
BEGIN
    SELECT RAISE(FAIL, 'Return date cannot be before loan date.');
END
DELIMITER ;;
