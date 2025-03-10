-- สร้างตาราง books
CREATE TABLE books (
  id SERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- สร้าง function สำหรับอัพเดท updated_at โดยอัตโนมัติ
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- สร้าง trigger เพื่อเรียกใช้ function update_modified_column
CREATE TRIGGER update_books_modtime
BEFORE UPDATE ON books
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

-- สร้าง index บน title เพื่อเพิ่มประสิทธิภาพการค้นหา
CREATE INDEX idx_books_title ON books(title);

-- เพิ่มข้อมูลตัวอย่าง
INSERT INTO books (title) VALUES 
    ('Fundamental of Deep Learning in Practice'),
    ('Practical DevOps and Cloud Engineering'),
    ('Mastering Golang for E-commerce Back End Development');