package auth

import(
	"time"
	"errors"

	"github.com/dgrijalva/jwt-go"
)


var secretKey = []byte("rest-api")

type Claims struct{
	Email string
	jwt.StandardClaims
}

func GenerateJWT(email string) string{
	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
		 ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Время истечения токена
		},
	   }
	  
	   // Подписание токена
	   token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	   tokenString, _ := token.SignedString(secretKey)
	   return tokenString
}

func ValidateToken(tokenString string) (*Claims, error) {
	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
	 return secretKey, nil
	})
	if err != nil {
	 return nil, err
	}
   
	// Проверка подписи
	if !token.Valid {
	 return nil, errors.New("неверный токен")
	}
   
	// Преобразование claims в MyClaims
	claims, ok := token.Claims.(*Claims)
	if !ok {
	 return nil, errors.New("неверный тип claims")
	}
   
	return claims, nil
}