package main

import "fmt"
import "encoding/json"


type Movie struct {
	Title	string	`json:"title"`
	Year	int		`json:"year"`
	Price	int
	Actors	[]string
}

func main()	{
	movie := Movie{"ET", 2000, 100, []string{"A", "b", "c"}}

	jsonStr, err := json.Marshal(movie)
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Printf("%s\n", jsonStr)
	
	my_movie := Movie{}
	err = json.Unmarshal(jsonStr, &my_movie)
	if err != nil {
		fmt.Println("err")
	}
	fmt.Printf("%v\n", my_movie)
}