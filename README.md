- **Структура `Database`:** Содержит соединение и логгер. Логгер используется для записи информации о выполнении запросов и ошибках.
- **Функция `NewDatabase`:** Создаёт новое соединение с базой данных и возвращает структуру `Database`.
- **Функция `Close`:** Закрывает соединение с базой данных.
- **Функция `Exec`:** Выполняет SQL-запросы, которые не возвращают результат (например, `INSERT`, `UPDATE`).
- **Функция `Query`:** Выполняет SQL-запросы, которые возвращают результат, и сканирует их в предоставленный объект.
- **Функция `Transaction`:** Обрабатывает выполнение операций в транзакции. В случае возникновения ошибки или паники, транзакция откатывается.

### Использование:

Подключение логгера и выполнение операций с базой данных:

```go
package main

import (
    "context"
    "log"
    "os"
    "mypgx"
)

func main() {
    logger := log.New(os.Stdout, "DB_LOG: ", log.LstdFlags)
    db, err := mypgx.NewDatabase("postgres://user:password@localhost:5432/dbname", logger)
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer db.Close()

    // Пример выполнения запроса
    _, err = db.Exec("CREATE TABLE example (id SERIAL PRIMARY KEY, name TEXT)")
    if err != nil {
        log.Fatalf("Failed to execute query: %v", err)
    }

    // Пример выполнения и сканирования запроса
    var results []struct {
        ID   int
        Name string
    }
    err = db.Query(&results, "SELECT id, name FROM example")
    if err != nil {
        log.Fatalf("Failed to query data: %v", err)
    }

    for _, result := range results {
        logger.Printf("Row: ID=%d, Name=%s", result.ID, result.Name)
    }

    // Пример транзакции
    err = db.Transaction(context.Background(), func(tx pgx.Tx) error {
        _, err := tx.Exec(context.Background(), "INSERT INTO example (name) VALUES ($1)", "John Doe")
        if err != nil {
            return err
        }
        // Другие операции в рамках транзакции
        return nil
    })

    if err != nil {
        log.Fatalf("Transaction failed: %v", err)
    }
}
