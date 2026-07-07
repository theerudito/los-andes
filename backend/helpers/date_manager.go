package helpers

import "time"

func FechaActual() string {
	return time.Now().Format("02/01/2006 15:04:05")
}
