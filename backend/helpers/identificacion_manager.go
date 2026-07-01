package helpers

func TipoIdentificacion(valor string) string {

	switch len(valor) {
	case 13:
		return "R"
	case 10:
		return "C"
	default:
		return "P"
	}
}
