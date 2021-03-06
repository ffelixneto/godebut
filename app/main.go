package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"environ"
)

func main() {
	// Cachacas = []Cachaca{
	// 	{Nome: "51", Volume: "974ml", Custo: "8"},
	// 	{Nome: "Matuta", Volume: "1000ml", Custo: "30"},
	// }

	// Endpoints = []Endpoint{
	// 	{},
	// }

	initialeMigration()
	fmt.Println("Conectando na base teste - OK ...")

	debut()
}

// pageInitial (Home Page) função para mostrar uma home page qualquer
// func pageInitial(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint: pageInitial")

	// fmt.Fprintln(w, "- ......... -")
	// fmt.Fprintln(w, "-- 2ez4Flx --")
	// fmt.Fprintln(w, "-  GoDrink  -")
	// fmt.Fprintln(w, "-  BzrSytes -")
	// fmt.Fprintln(w, "- ......... -")
// }

// aProposDe (About) é para mostrar minhas informações
func aProposDe(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: aproposde")

	qui := "Felix Neto"

	fmt.Fprint(w, "A propos de...", qui)
}

// toutesv1Endpoints é o endpoint para exibir listagem de endpoints v1
// func toutesv1Endpoints(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint: toutesv1Endpoints")
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(Endpoints)
// }

// toutesCachacas é o endpoint para listar todas as cachaças cadastradas
func toutesCachacas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: toutesCachacas")
	w.Header().Set("Content-Type", "application/json")
	db := PostgresConn()
	defer db.Close()

	fmt.Println("Listando cachacas...")
	var cachacas []Cachaca
	db.Find(&cachacas)

	for _, cana := range cachacas {
		fmt.Println("Id", cana.ID, "Nome: ", cana.Nome, "Volume: ", cana.Volume, "Custo: ", cana.Custo)
	}
	json.NewEncoder(w).Encode(cachacas)
}

// uneCachaca é o endpoint para listar uma cachaça buscando pelo nome informado na URL
func uneCachaca(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: uneCachaca")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]

	var cachaca Cachaca

	if cle != "" {
		db.Find(&cachaca, cle)

		if cachaca.ID != 0 {
			fmt.Println("Selecionando cachaca com Id:", cle)
			json.NewEncoder(w).Encode(cachaca)
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhuma cacaça econtrada com o Id:", cle)
			fmt.Fprintln(w, "Nenhuma cacaça econtrada com o Id:", cle)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("O dado de ID deve ser informado !")
		fmt.Fprintln(w, "O dado de ID deve ser informado !")
	}
}

// nouvelleCachaca é o endpoint para criar novos registros de cachaca
func nouvelleCachaca(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: nouvelleCachaca")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	reqBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		panic(erro)
	}
	var cachaca Cachaca
	json.Unmarshal(reqBody, &cachaca)

	// Validar dados basicos antes de adicionar
	if cachaca.Nome != "" && cachaca.Volume != "" && cachaca.Custo != "" {
		db.Create(&Cachaca{Nome : cachaca.Nome, Volume: cachaca.Volume, Custo: cachaca.Custo})
		fmt.Fprintln(w, "A cachaca ", cachaca.Nome, " foi adicionada a lista.")
		json.NewEncoder(w).Encode(cachaca)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Os dados de \"nome\", \"volume\" e \"custo\" devem ser informados !")
	}
}

// renouvelleCachaca é o endpoint para atualizar um registro de cachaca
func renouvelleCachaca(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: renouvelerCachaca")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]
	reqBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		panic(erro)
	}

	var new_cachaca Cachaca
	var cachaca Cachaca
	json.Unmarshal(reqBody, &new_cachaca)

	if cle != "" {
		db.Find(&cachaca, cle)

		if cachaca.ID != 0 {
			fmt.Println(
				"Dados encontrados", 
				cachaca.ID,
				cachaca.Nome, 
				cachaca.Volume,
				cachaca.Custo,
			)

			db.Model(&cachaca).Updates(
				Cachaca{
					Nome: new_cachaca.Nome,	
					Volume : new_cachaca.Volume, 
					Custo : new_cachaca.Custo,
				})

			db.Find(&cachaca, cle)
			fmt.Println(
				"Dados atualizados", 
				cachaca.ID,
				cachaca.Nome, 
				cachaca.Volume, 
				cachaca.Custo, 
			)
			json.NewEncoder(w).Encode(cachaca)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhuma cachaca encontrada com o Id:", cle)
			fmt.Fprintln(w, "Nenhuma cachaca encontrada com o Id:", cle)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "O dado de ID deve ser preenchido !")
		fmt.Println("O dado de ID deve ser preenchido !")
	}
}

// effacerCachaca é o endpoint para deletar
func effacerCachaca(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: effacerCachaca")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]
	
	var cachaca Cachaca

	if cle != "" {
		db.Find(&cachaca, cle)

		if cachaca.ID != 0 {
			fmt.Println("Consumidor ID", cachaca.ID, "deletada ! ")
			fmt.Fprintln(w, "Consumidor ID", cachaca.ID, "deletada ! ")
			db.Delete(&cachaca)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhuma cachaca encontrada com o Id:", cle )
			fmt.Fprintln(w, "Nenhuma cachaca encontrada com o Id:", cle)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("O dado de ID deve ser preenchido !")
		fmt.Fprintln(w, "O dado de ID deve ser preenchido !")
	}
}

// toutesConsumidores é o endpoint para listar todos os consumidores
func toutesConsumidores(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: toutesConsumidores")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	fmt.Println("Listando consumidores...")
	var consumidores []Consumidor
	db.Find(&consumidores)

	for _, consu := range consumidores {
		fmt.Println("Id", consu.ID, "Nome: ", consu.Nome, "Idade: ", consu.Idade)
	}
	json.NewEncoder(w).Encode(consumidores)
}

// uneConsumidor é o endpoint para listar um consumidor buscando pelo nome informado na URL
func uneConsumidor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: uneConsumidor")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]

	var consumidor Consumidor

	if cle != "" {
		db.Find(&consumidor, cle)

		if consumidor.ID != 0 {
			fmt.Println("Selecionado consumidor com Id", cle)
			json.NewEncoder(w).Encode(consumidor)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhum consumidor encontrado com o Id:", cle)
			fmt.Fprintln(w, "Nenhum consumidor encontrado com o Id:", cle)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("O dado de ID deve ser informado !")
		fmt.Fprintln(w, "O dado de ID deve ser informado !")
	}
}

// nouvelleConsumidor é o endpoint para adicionar um novo registro de consumidor
func nouvelleConsumidor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: nouvelleConsumidor")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	reqBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		panic(erro)
	}
	var consumidor Consumidor
	json.Unmarshal(reqBody, &consumidor)

	// Validar dados basicos antes de adicionar
	if consumidor.Nome != "" && consumidor.Idade != "" {
		db.Create(&Consumidor{Nome: consumidor.Nome, Idade: consumidor.Idade})
		fmt.Fprintln(w, "O consumidor ", consumidor.Nome, " foi adicionado a lista.")
		json.NewEncoder(w).Encode(consumidor)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Os dados de \"nome\" e \"idade\" devem ser informados !")
	}
}

// renouvelleConsumidor é o endpoint para atualizar um registro de consumidor
func renouvelleConsumidor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: renouvelleConsumidor")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]
	reqBody, erro := ioutil.ReadAll(r.Body)
	if erro != nil {
		panic(erro)
	}

	var new_consumidor Consumidor
	var consumidor Consumidor
	json.Unmarshal(reqBody, &new_consumidor)

	if cle != "" {
		db.Find(&consumidor, cle)

		if consumidor.ID != 0 {
			fmt.Println(
				"Dados encontrador", 
				consumidor.ID,
				consumidor.Nome, 
				consumidor.Idade,
			)

			db.Model(&consumidor).Updates(
				Consumidor{
					Nome: new_consumidor.Nome,	
					Idade : new_consumidor.Idade,
				})
			
			db.Find(&consumidor, cle)
			fmt.Println(
				"Dados atualizados", 
				consumidor.ID,
				consumidor.Nome,
				consumidor.Idade, 
			)
			json.NewEncoder(w).Encode(consumidor)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhum consumidor encontrado com o Id:", "\""+cle+"\"")
			fmt.Fprintln(w, "Nenhum consumidor encontrado com o Id: \""+cle+"\"")
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("O dado de ID deve ser preenchido !")
		fmt.Fprintln(w, "O dado de ID deve ser preenchido !")
	}
}

// effacerConsumidor é o endpoint para deletar um registro de consumidor
func effacerConsumidor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: effacerConsumidor")
	w.Header().Set("Content-Type", "application/json")

	db := PostgresConn()
	defer db.Close()

	vars := mux.Vars(r)
	cle := vars["id"]

	var consumidor Consumidor

	if cle != "" {
		db.Find(&consumidor, cle)

		if consumidor.ID != 0 {
			fmt.Println("Consumidor ID", consumidor.ID, "deletado ! ")
			fmt.Fprintln(w, "Consumidor ID", consumidor.ID, "deletado ! ")
			db.Delete(&consumidor)

		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Nenhum consumidor encontrado com o Id:", cle)
			fmt.Fprintln(w, "Nenhum consumidor encontrado com o Id:", cle)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("O dado de ID deve ser preenchido !")
		fmt.Fprintln(w, "O dado de ID deve ser preenchido !")
	}
}

// initialeMigration é a função para executar a migração no banco
// e fazer o seed com alguns valores para teste
// func initialeMigration(db *gorm.DB) {
func initialeMigration() {

	db := PostgresConn()

	// Limpeza de dados anteriores
	db.DropTable(&Consumidor{})
	db.AutoMigrate(&Consumidor{})
	db.DropTable(&Cachaca{})
	db.AutoMigrate(&Cachaca{})

	// Inserção de dados de cachaças para teste
	db.Create(&Cachaca{Nome: "51", Volume: "974ml", Custo: "10"})
	db.Create(&Cachaca{Nome: "Matuta", Volume: "1000ml", Custo: "30"})
	db.Create(&Cachaca{Nome: "Carangueijo", Volume: "1000ml", Custo: "25"})

	// Inserção de dados de consumidores para teste
	db.Create(&Consumidor{Nome: "felix marmotinha", Idade: "18"})
	db.Create(&Consumidor{Nome: "felix-2-devops", Idade: "33"})
	db.Create(&Consumidor{Nome: "jorbson-2-scripts", Idade: "21"})
	db.Create(&Consumidor{Nome: "paulo-2-manager", Idade: "28"})
	db.Create(&Consumidor{Nome: "Nicer", Idade: "1"})
	db.Create(&Consumidor{Nome: "Ezy", Idade: "2"})
	db.Delete(&Consumidor{}, 1)
}

// PostgresConn tenta conectar e retornar uma conexão ao banco de dados postgres com retry
// Se o numero de retry esgotar, gera panic(err)
func PostgresConn() *gorm.DB {
	var db *gorm.DB
	var err error

	db_host := environ.GetEnvironValue("DB_HOST")
	db_user := environ.GetEnvironValue("DB_USER")
	db_passwd := environ.GetEnvironValue("DB_PASSWD")
	db_name := environ.GetEnvironValue("DB_NAME")
	db_port := environ.GetEnvironValue("DB_PORT")

	db_try := 40
	var db_url = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		db_host,
		db_user,
		db_passwd,
		db_name,
		db_port)

	db, err = gorm.Open("postgres", db_url)
	for err != nil {
		fmt.Println("Tentativa conexão ao banco...", db_try)

		if db_try > 1 {
			db_try--
			time.Sleep(5 * time.Second)
			db, err = gorm.Open("postgres", db_url)
			continue
		}
		panic(err)
	}
	return db
}

// debut é a função que vai ativar as rotas da API
func debut() {

	roteur := mux.NewRouter().StrictSlash(true)

	// Rotas de visualização
	// roteur.HandleFunc("/", pageInitial)
	
	roteur.HandleFunc("/aproposde", aProposDe)
	
	// Rotas de documentação
	// roteur.HandleFunc("/v1/toutesendpoints", toutesv1Endpoints)
	// TODO - ENDPOINT PARA LSITAGEM DE ENDPOINTS V1
	
	// Rotas de endpoints das cachaças
	roteur.HandleFunc("/v1/toutescachacas", toutesCachacas)
	roteur.HandleFunc("/v1/unecachaca", nouvelleCachaca).Methods("POST")
	roteur.HandleFunc("/v1/unecachaca/{id}", renouvelleCachaca).Methods("PUT")
	roteur.HandleFunc("/v1/unecachaca/{id}", effacerCachaca).Methods("DELETE")
	roteur.HandleFunc("/v1/unecachaca/{id}", uneCachaca)
	
	// Rotas de endpoints dos consmidores
	roteur.HandleFunc("/v1/toutesconsumidores", toutesConsumidores)
	roteur.HandleFunc("/v1/uneconsumidor", nouvelleConsumidor).Methods("POST")
	roteur.HandleFunc("/v1/uneconsumidor/{id}", renouvelleConsumidor).Methods("PUT")
	roteur.HandleFunc("/v1/uneconsumidor/{id}", effacerConsumidor).Methods("DELETE")
	roteur.HandleFunc("/v1/uneconsumidor/{id}", uneConsumidor)
	
	// Pagina inicial
	roteur.PathPrefix("/").Handler(http.StripPrefix("/",http.FileServer(http.Dir("./pages/"))))
	
	api_environ := environ.GetEnvironValue("API_PORT")
	var api_port = ":" + fmt.Sprintf("%v", api_environ)
	
	fmt.Println("Iniciando alambique em http://localhost" + api_port)
	log.Fatal(http.ListenAndServe(api_port, roteur))
}

// Cachaca é a estrutura base para a tabela "cachacas" no banco
type Cachaca struct {
	gorm.Model
	Nome   string
	Volume string
	Custo  string
}

// Consumidor é a estrutura base para a tabela "consumidors" no banco
type Consumidor struct {
	gorm.Model
	Nome  string
	Idade string
}

// Endpoint é a estrutura base para o projeto cachaça
// type Endpoint struct {
// 	gorm.Model
// 	Nome   	string
// 	Versao 	string
// 	Metodo  string
// }

var Cachacas []Cachaca
