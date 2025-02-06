-- Comment
CREATE TABLE t4 (
    -- in SQL body comment
    name TEXT
);

-- Another comment
-- Multiple comments
/* Other type of comment */

CREATE TRIGGER track_inserts
AFTER INSERT ON t1
FOR EACH ROW
BEGIN
  INSERT INTO t4 (name)

  VALUES ("new");
END;
