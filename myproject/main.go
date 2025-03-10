package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
    "context"
    "log"
    "time"
    // "myproject/internal/bookstore"
)


type BookDatabase interface {
    GetBook(id int) (string, error)
    AddBook(title string) error
    DeleteBook(id int) error
	GetAllBooks() ([]string, error)
    Close() error
}

type PostgresDatabase struct {
    db *sql.DB
}

func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %v", err)
    }

    // ตั้งค่า connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)
    db.SetConnMaxLifetime(5 * time.Minute)

    // ทดสอบการเชื่อมต่อด้วย context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }

    return &PostgresDatabase{db: db}, nil
}

func (pdb *PostgresDatabase) GetBook(ctx context.Context, id int) (string, error) {
    var title string
    err := pdb.db.QueryRowContext(ctx, "SELECT title FROM books WHERE id = $1", id).Scan(&title)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", fmt.Errorf("book not found")
        }
        return "", fmt.Errorf("failed to get book: %v", err)
    }
    return title, nil
}

func (pdb *PostgresDatabase) GetAllBooks() ([]string, error) {
	rows, err := pdb.db.Query("SELECT title FROM books")
	if err != nil {
		return nil, fmt.Errorf("failed to query books: %v", err)
	}
	defer rows.Close()

	var titles []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, fmt.Errorf("failed to scan book title: %v", err)
		}
		titles = append(titles, title)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over books: %v", err)
	}

	return titles, nil
}

func (pdb *PostgresDatabase) AddBook(ctx context.Context, title string) error {
    _, err := pdb.db.ExecContext(ctx, "INSERT INTO books (title) VALUES ($1)", title)
    if err != nil {
        return fmt.Errorf("failed to add book: %v", err)
    }
    return nil
}

func (pdb *PostgresDatabase) DeleteBook(ctx context.Context, id int) error {
    result, err := pdb.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
    if err != nil {
        return fmt.Errorf("failed to delete book: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("book not found")
    }

    return nil
}

func (pdb *PostgresDatabase) Close() error {
	return pdb.db.Close()
}



func main() {
    // สร้างการเชื่อมต่อกับ Database
    db, err := NewPostgresDatabase("host=localhost port=5432 user=bookstore_user password=your_strong_password dbname=bookstore sslmode=disable")
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // สร้าง BookStore
    // store := bookstore.NewBookStore(db)

    // สร้าง Context พร้อม Timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // เพิ่มหนังสือ
    err = db.AddBook(ctx, "The Go Programming Language")
    if err != nil {
        log.Printf("Failed to add book: %v", err)
    } else {
        fmt.Println("Book added successfully")
    }

    // ดึงข้อมูลหนังสือ
    title, err := db.GetBook(ctx, 1) // สมมติว่า id = 1
    if err != nil {
        log.Printf("Failed to get book: %v", err)
    } else {
        fmt.Printf("Book title: %s\n", title)
    }

    // ลบหนังสือ
    err = db.DeleteBook(ctx, 1) // สมมติว่า id = 1
    if err != nil {
        log.Printf("Failed to delete book: %v", err)
    } else {
        fmt.Println("Book deleted successfully")
    }
}