package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/Bloxico/exchange-gateway/sofija/server/params"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/julienschmidt/httprouter"
)

// The signing key used to create JWT
var (
	jwtSigningKey = []byte("DebugSigningKey")
)

// Where the user's email and ID will be stored in the request context
const USER_EMAIL_CTX_KEY = "EMAIL"
const USER_ID_CTX_KEY = "ID"

type CustomClaims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.RegisteredClaims
}

// Creates a JWT for the given email
// Returns the JWT, or an error
func CreateJWT(email string, id string) (string, error) {

	claims := &CustomClaims{
		email,
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "EGW",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	ss, err := token.SignedString(jwtSigningKey)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// utility method: creates a JWT auth header in the expected format (`Bearer xxx,yyy,xxx`)
func WrapJWTHeader(token string) string {
	return "Bearer " + token
}

// unwraps the auth header, and returns the JWT token, in form of xxx.yyy.zzz
func UnwrapJWTHeader(authHeader string) string {
	return strings.Split(authHeader, " ")[1]
}

// Checks whether the given auth header is valid
func IsJWTHeaderValid(authHeader string) bool {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	jwtToken := strings.Split(authHeader, " ")
	if len(jwtToken) < 2 {
		return false
	}
	// since the JWT format is xxx.yyy.zzz, we need to split by .
	parts := strings.Split(jwtToken[1], ".")
	// payload is composed of xxx.yyy, and the signature is zzz
	payload := strings.Join(parts[0:2], ".")
	err := jwt.SigningMethodHS512.Verify(payload, parts[2], jwtSigningKey)

	return err == nil
}

// Extracts the claims from the given JWT
func GetJWTClaims(token string) (jwt.MapClaims, error) {
	res, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) { return jwtSigningKey, nil })
	if err != nil {
		return nil, err
	}
	if claims, ok := res.Claims.(jwt.MapClaims); ok && res.Valid {
		return claims, nil
	} else {
		return nil, errors.New("error parsing JWT claims")
	}
}

// Authenticates requests by checking the JWT
// If authentication is successful, the user's email will be attached to the request context
func AuthJWT(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authHeader := req.HeaderParameter("Authorization")

	if !IsJWTHeaderValid(authHeader) {
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}

	// unpack JWT
	claims, err := GetJWTClaims(UnwrapJWTHeader(authHeader))
	if err != nil {
		resp.WriteErrorString(500, "Server error")
	}

	// attach user email to request context
	userEmail := claims["email"].(string)
	userId := claims["id"].(string)
	updated := params.WithRequest(req.Request, httprouter.Params{
		httprouter.Param{Key: USER_EMAIL_CTX_KEY, Value: userEmail},
		httprouter.Param{Key: USER_ID_CTX_KEY, Value: userId},
	})
	req.Request = updated

	chain.ProcessFilter(req, resp)
}
