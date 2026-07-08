package helpers

import (
	"database/sql"
	"fmt"
)

func ObtenerCodigo(exec ExecutorDB, prefijo string) (string, error) {

	var actual, digitos int

	query := `SELECT actual, digitos FROM secuencial WHERE prefijo = ?`
	err := exec.QueryRow(query, prefijo).Scan(&actual, &digitos)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("el prefijo '%s' no existe en la tabla secuencial", prefijo)
		}
		return "", err
	}

	var codigoFormateado string

	if prefijo == "O" {
		formato := fmt.Sprintf("%%0%dd", digitos)
		codigoFormateado = fmt.Sprintf(formato, actual)
	} else {
		formato := fmt.Sprintf("%s%%0%dd", prefijo, digitos-1)
		codigoFormateado = fmt.Sprintf(formato, actual)
	}

	return codigoFormateado, nil
}

func ActualizarCodigo(exec ExecutorDB, prefijo string) error {

	query := `
    UPDATE secuencial 
    SET actual = actual + 1 
    WHERE prefijo = ?`

	result, err := exec.Exec(query, prefijo)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no se pudo actualizar: el prefijo '%s' no existe", prefijo)
	}

	return nil
}
