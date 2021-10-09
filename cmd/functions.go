package cmd

import (
	"creativeadvtech/pkg"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func InitRedis() {
	var b pkg.BodyJson
	b.Flush()

	json, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}
	err = Conn.Set("foo", json, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func GetValRedisString() float64 {
	var answer pkg.BodyJson
	val, err := Conn.Get("foo").Result()
	if err != nil {
		InitRedis()
	}
	err = json.Unmarshal([]byte(val), &answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
	}
	return float64(answer.Count)
}

//func CheckError(err error, w http.ResponseWriter){
//	if err != nil {
//		logrus.Warnln(http.StatusInternalServerError)
//		w.WriteHeader(http.StatusInternalServerError)
//		fmt.Fprintf(w, "Something wrong %s", err)
//	}
//}
