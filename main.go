package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"math/rand"
	"net/http"
	"strconv"
)

type Pharmacy struct {
	Id        int    `json:"id"`
	Client_id int    `json:"client_id"`
	Name      string `json:"name"`
}

type GenericError struct {
	Msg string `json:"msg"`
	Ids []int  `json:"ids"`
}

func NewGenericError(msg string, ids []int) *GenericError {
	return &GenericError{
		Msg: msg,
		Ids: ids,
	}
}

func Get_pharm_by_id(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, _ := strconv.Atoi(ps.ByName("id"))

	// this is an example of the resource not being found
	if id == -1 {
		w.WriteHeader(404)
		return
	}

	m := Pharmacy{
		Id:        id,
		Client_id: 1,
		Name:      fmt.Sprintf("Pharmacy_%d", id),
	}

	j, _ := json.Marshal(m)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", j)
}

func Add_update_pharmacy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := Pharmacy{}
	json.NewDecoder(r.Body).Decode(&p)
	p.Id = rand.New(rand.NewSource(99)).Int()

	// how we can handle an error
	if p.Client_id == 0 {
		w.WriteHeader(400)
		j, _ := json.Marshal(NewGenericError("client_id is invalid", nil))
		fmt.Fprintf(w, "%s", j)
		return
	}

	// how we can return that we created the new resource
	// Returning the location header is a "best practice", it allows
	// the client to extract that header and then get back the created resource
	w.Header().Set("Location", fmt.Sprintf("/pharmacies/%d", p.Id))
	// 201: Created
	w.WriteHeader(201)
}

func Get_client_pharmacies(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client_id, _ := strconv.Atoi(ps.ByName("client_id"))
	x := 100
	pr := make([]Pharmacy, x)
	for i := 0; i < x; i++ {
		pr[i] = Pharmacy{
			Id:        i,
			Client_id: client_id,
			Name:      fmt.Sprintf("Pharmacy_%d", i),
		}
	}

	j, _ := json.Marshal(pr)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", j)
}

func main() {
	router := httprouter.New()
	router.GET("/pharmacies/:id", Get_pharm_by_id)
	router.POST("/pharmacies", Add_update_pharmacy)

	router.GET("/client_id/:client_id/pharmacies", Get_client_pharmacies)

	http.ListenAndServe("localhost:8080", router)
}
