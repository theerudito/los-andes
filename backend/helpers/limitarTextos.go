package helpers

func Limitar(s string, limite int) string {
	if len(s) > limite {
		return s[:limite-3] + "..."
	}
	return s
}
