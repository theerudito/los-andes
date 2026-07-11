package helpers

import "strings"

func ObtenerPalabra(nombres, apellidos string) string {

	n := strings.Fields(nombres)
	a := strings.Fields(apellidos)

	var nombre string
	var apellido string

	if len(n) > 0 {
		nombre = n[0]
	}

	if len(a) > 0 {
		apellido = a[0]
	}

	if nombre != "" && apellido != "" {
		return nombre + " " + apellido
	}

	if nombre != "" {
		return nombre
	}

	if apellido != "" {
		return apellido
	}

	return ""
}
