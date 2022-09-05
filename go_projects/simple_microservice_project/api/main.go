package main

import(
	"fmt"
	"net/http"
	"os"
	"log"
	"github.com/dgrijalva/jwt-go"
)
var MySigningKey = []byte(os.Getenv("SECRET_KEY_API"))

func homePage(w  http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Super Secret Information")
}
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if r.Header["Token"] != nil{
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token)(interface{}, error){
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
					return nil, fmt.Errorf("There was an error")
				}
				aud:= "billing.jwtgo.io"
				checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
				if !checkAudience{
					return nil, fmt.Errorf("Invalid audience.")
				return MySigningKey, nil
			}
			iss:= "jwtgo.io"

			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss{
				return nil, fmt.Errorf("Invalid issuer.")
			}
			return MySigningKey, nil
		})
		if err != nil{
			fmt.Fprintf(w, err.Error())
			
		}
		if token.Valid{
			endpoint(w, r)
		}
	}else{
			fmt.Fprintf(w, "Not Authorized token valid")
		}
	})
}


func handleRequests(){
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main(){
	fmt.Println("API Server Started")
	handleRequests()
}