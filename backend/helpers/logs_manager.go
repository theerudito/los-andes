package helpers

func InsertLogsError(exec ExecutorDB, tableName string, message string) error {
	query := `
        INSERT INTO log_error (modulo, mensaje_error)
        VALUES (?, ?)`

	_, err := exec.Exec(query, tableName, message)
	return err
}

func InsertLogs(exec ExecutorDB, action string, tableName string, recordId int, description string) error {
	query := `
        INSERT INTO log_ok (modulo, accion, registro_id, descripcion)
        VALUES (?, ?, ?, ?)`

	_, err := exec.Exec(query, tableName, action, recordId, description)
	return err
}
