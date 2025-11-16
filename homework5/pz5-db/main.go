package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// fallback — прямой DSN в коде (только для учебного стенда!)
	dsn := "postgres://user:password@localhost:5432/todo?sslmode=disable"

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты"}
	for _, title := range titles {
		id, err := repo.CreateTask(ctx, title)
		if err != nil {
			log.Fatalf("CreateTask error: %v", err)
		}
		log.Printf("Inserted task id=%d (%s)", id, title)
	}

	// 2) Прочитаем список задач
	tasks, err := repo.ListTasks(ctx)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем все задачи
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// 4) Обновляем таску под id=1
	err = repo.MarkTaskDone(ctx, 1)
	if err != nil {
		log.Fatalf("Not found task id=1: %v", err)
	}

	// 5) Выводим все таски done=True
	tasksDone, err := repo.ListDone(ctx, true)
	if err != nil {
		log.Fatalf("ListTasksDone error: %v", err)
	}
	fmt.Println("=== Tasks Done True===")
	for _, t := range tasksDone {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// 6) Выводим таску id=1
	t, err := repo.FindByID(ctx, 1)
	if err != nil {
		log.Fatalf("ListTasksDone error: %v", err)
	}
	fmt.Println("=== Tasks ID=1===")

	fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
		t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))

	// 7) Вставим задачи новым методом
	titles = []string{"test1", "test2"}
	err = repo.CreateMany(ctx, titles)
	if err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}

	// 8) Напечатаем все задачи
	tasks, err = repo.ListTasks(ctx)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}
	fmt.Println("===  + new tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

}
