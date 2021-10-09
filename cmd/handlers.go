package cmd

import (
	"creativeadvtech/pkg"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	opsProcessed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "Counter",
			Help: "Shows the value of Counter",
		},
		[]string{"custom_label"},
	)
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.With(prometheus.Labels{"custom_label": "my"}).Set(GetValRedisString())
			time.Sleep(10 * time.Second)
		}
	}()
}

func Health(w http.ResponseWriter, r *http.Request) {
	pong, err := Conn.Ping().Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("FAILED connect to Redis"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(pong))
	}
}

//func Metrics(w http.ResponseWriter, r *http.Request) {
//	w.WriteHeader(http.StatusOK)
//}

func GetCounter(w http.ResponseWriter, r *http.Request) {
	var answer pkg.BodyJson

	val, err := Conn.Get("foo").Result()
	if err != nil {
		InitRedis()
	}

	err = json.Unmarshal([]byte(val), &answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(answer)
}

func IncrementCounter(w http.ResponseWriter, r *http.Request) {
	var answer pkg.BodyJson

	val, err := Conn.Get("foo").Result()
	if err != nil {
		InitRedis()
	}

	err = json.Unmarshal([]byte(val), &answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	answer.Count++

	jsonData, err := json.Marshal(answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	err = Conn.Set("foo", jsonData, 0).Err()
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(answer)
}

func UpdateCounter(w http.ResponseWriter, r *http.Request) {
	var answer pkg.BodyJson

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	err = json.Unmarshal(body, &answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	jsonData, err := json.Marshal(answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	err = Conn.Set("foo", jsonData, 0).Err()
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(answer)
}

func DeleteCounter(w http.ResponseWriter, r *http.Request) {
	InitRedis()

	var answer pkg.BodyJson

	val, err := Conn.Get("foo").Result()
	if err != nil {
		InitRedis()
	}

	err = json.Unmarshal([]byte(val), &answer)
	if err != nil {
		logrus.Warnln(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Something wrong %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(answer)
}
