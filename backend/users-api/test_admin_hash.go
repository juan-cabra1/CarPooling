package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password del admin
	password := "admin123!"

	// Hash que está en la BD (bcrypt cost 12)
	hashInDB := "$2a$12$xpY6Qi0yXnmn2eAWI/2I7eppTTvPEu/wSA0Tq/oL.yFpTa8U6K6Qe"

	// Verificar si el hash coincide
	err := bcrypt.CompareHashAndPassword([]byte(hashInDB), []byte(password))
	if err != nil {
		fmt.Printf("❌ El hash NO coincide con la contraseña '%s'\n", password)
		fmt.Printf("Error: %v\n", err)

		// Generar un nuevo hash correcto
		fmt.Println("\nGenerando nuevo hash con bcrypt cost 10...")
		newHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			fmt.Printf("Error generando hash: %v\n", err)
			return
		}
		fmt.Printf("✅ Nuevo hash:\n%s\n", string(newHash))
	} else {
		fmt.Printf("✅ El hash coincide perfectamente con la contraseña '%s'\n", password)
	}
}
