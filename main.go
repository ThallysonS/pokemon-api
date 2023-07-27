package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/lib/pq"
)

type Pokemon struct {
	Nome   string `json:"name"`
	Altura int    `json:"height"`
	Peso   int    `json:"weight"`
}

func PegarDados(id int) (*Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d", id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var pokemon Pokemon
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return nil, err
	}
	return &pokemon, nil
}

func cadastrarNoBanco(pokemon *Pokemon, db *sql.DB) error {
	// Preparar a instrução SQL para inserção
	stmt, err := db.Prepare("INSERT INTO pokemons(nome, altura, peso) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Executar a instrução SQL com os valores do Pokémon
	_, err = stmt.Exec(pokemon.Nome, pokemon.Altura, pokemon.Peso)
	return err
}

func main() {
	// Conectar ao banco de dados PostgreSQL
	db, err := sql.Open("postgres", "user=postgres password=4002 dbname=pokemons sslmode=disable")
	if err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %s\n", err)

	}
	defer db.Close()

	PokeID := 26
	pokemon, err := PegarDados(PokeID)
	if err != nil {
		fmt.Printf("Erro ao obter os dados do Pokémon: %s\n", err)
		return
	}

	// Cadastrar o Pokémon no banco de dados
	err = cadastrarNoBanco(pokemon, db)
	if err != nil {
		fmt.Printf("Erro ao cadastrar o Pokémon no banco de dados: %s\n", err)
		return
	}

	fmt.Println("Dados do Pokémon cadastrados com sucesso!")
}
