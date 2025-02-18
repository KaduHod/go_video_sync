package utils
import "golang.org/x/crypto/bcrypt"
func HashPassword(password string) (string, error) {
	// Gera o hash usando o custo padrão (10)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Função para verificar se uma senha corresponde ao hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
