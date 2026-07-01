package helpers

func InsertLogsError(exec ExecutorDB, tableName string, message string) error {
	query := `
        INSERT INTO logs_error (table_name, message)
        VALUES ($1, $2)`
	_, err := exec.Exec(query, tableName, message)
	return err
}

func InsertLogs(exec ExecutorDB, action string, tableName string, recordId int, description string) error {
	query := `
        INSERT INTO log_actions (action, table_name, record_id, description)
        VALUES ($1, $2, $3, $4)`
	_, err := exec.Exec(query, action, tableName, recordId, description)
	return err
}
