package main
import (
	"fmt"
	"net/http"
	"os"
	"log"
	"time"
	jwt "github.com/dgrijalva/jwt-go"
)

var MySigningKey = []byte(os.Getenv("SECRET_KEY_JWTCREATOR"))

func GetJWT() (string, error){
	token:= jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true //info that you want to add
	claims["client"] = "Melike"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() //token expires in 30 minutes
	claims["aud"] = "billing.jwtgo.io"
	claims["iss"] = "jwtgo.io"
	tokenString, err:= token.SignedString(MySigningKey)
	if err!= nil{
		fmt.Errorf("something went wrong: %s", err.Error())
		return "", err
	}
return tokenString, nil

}
func Index(w http.ResponseWriter, r *http.Request){
	validToken, err := GetJWT()
	fmt.Println(validToken)
	if err!= nil{
		fmt.Println("Failed to generate token")
	}
	fmt.Fprintf(w, string(validToken))
}

func handleRequests(){
	http.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main(){
	handleRequests()
}