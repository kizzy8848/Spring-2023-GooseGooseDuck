package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
	"toj/define"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.StandardClaims
}

// GetMd5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("toj-key")

// GenerateToken
func GenerateToken(identity, name string, isAdmin int) (string, error) {
	UserClaim := &UserClaims{
		Identity:       identity,
		Name:           name,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendCode
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Get <tojproject@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "The verification code has been sent successfully, please check!"
	e.HTML = []byte("Your TOJ Verification Code(Valid for 10 Minutes): <b>" + code + "</b>")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "tojproject@163.com", "HFTHAYCSPVUMZHEI", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})

}

// GetUUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// GetRand
func GetRand() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

// CodeSave
func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/main.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	f.Write(code)
	defer f.Close()
	return path, nil
}

// CheckGoCodeValid
func CheckGoCodeValid(path string) (bool, error) {
	b, err := ioutil.ReadFile(path)

	fmt.Println("code")
	fmt.Println(string(b))

	if err != nil {
		return false, err
	}
	code := string(b)
	if len(code) == 0 {
		return false, nil
	}
	x := false
	for i := 0; i < len(code)-6; i++ {
		if code[i:i+6] == "import" {
			x = true
			var flag byte
			for i = i + 7; i < len(code); i++ {
				if code[i] == ' ' {
					continue
				}
				flag = code[i]
				break
			}
			if flag == '(' {
				for i = i + 1; i < len(code); i++ {
					if code[i] == ')' {
						break
					}
					if code[i] == '"' {
						t := ""
						for i = i + 1; i < len(code); i++ {
							if code[i] == '"' {
								break
							}
							t += string(code[i])
						}
						if _, ok := define.ValidGolangPackageMap[t]; !ok {
							return false, nil
						}
					}
				}
			} else if flag == '"' {
				t := ""
				for i = i + 1; i < len(code); i++ {
					if code[i] == '"' {
						break
					}
					t += string(code[i])
				}
				if _, ok := define.ValidGolangPackageMap[t]; !ok {
					return false, nil
				}
			}
		}

	}
	if !x{
		return false, nil
	}
	return true, nil
}
