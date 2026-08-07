package csv_handler

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/priyanshu-t-singh/dreamsofcode-goprojects/01-todo-list/internal/model"
)

type CSVHandler struct{}

func (h *CSVHandler) WriteToCSV(filename string, tasks []model.Task) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create a file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"ID", "Description", "CreatedAt", "IsCompleted"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, task := range tasks {
		record := []string{task.ID, task.Description, task.CreatedAt, strconv.FormatBool(task.IsCompleted)}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

func (h *CSVHandler) ReadFromCSV(filename string) ([]model.Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open a file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var tasks []model.Task
	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf("failed to read records: %w", err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		IsCompleted, err := strconv.ParseBool(record[3])
		if err != nil {
			return nil, fmt.Errorf("unable to parse IsCompleted %w", err)
		}

		task := model.Task{
			ID:          record[0],
			Description: record[1],
			CreatedAt:   record[2],
			IsCompleted: IsCompleted,
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
