package db

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

const (
	plantsTableName = "plants"
)

var allTableNames = []string{
	plantsTableName,
}

func (i *Instance) ensureSchemas(ctx context.Context, suffix string) error {
	// plants table
	pt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	    id STRING NOT NULL PRIMARY KEY,
	    name STRING NOT NULL,
	    botanical_name STRING NOT NULL,
	    status STRING,
        water_pref STRING,
        light_pref STRING,
        humidity_pref STRING
	)
	    `, plantsTableName)

	_, err := i.db.Exec(ctx, pt)
	if err != nil {
		return err
	}
	plants := generatePlants()
	pi := fmt.Sprintf(`
        INSERT INTO %s (id, name, botanical_name, status, water_pref, light_pref, humidity_pref)
        VALUES ($1,$2,$3,$4,$5,$6,$7)
    `, plantsTableName)
	for j := 0; j < len(plants); j++ {
		_, err := i.db.Exec(ctx, pi, plants[j].ID, plants[j].Name, plants[j].BotanicalName, plants[j].Status, plants[j].WaterPref, plants[j].LightPref, plants[j].HumidityPref)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Instance) cleanupSchemas() error {
	for _, table := range allTableNames {
		q := fmt.Sprintf("DROP TABLE %s;", table)
		fmt.Printf("Dropping table: %s\n", table)
		_, err := i.db.Exec(context.Background(), q)
		if err != nil {
			return err
		}

	}
	return nil
}

// max and min number or records
var (
	min        = 5
	max        = 30
	adjectives = []string{
		"pretty", "large", "big", "small", "tall", "short", "long",
		"leafy", "adapted", "weary", "stronger", "adorable", "inexpensive",
		"adaptable", "expensive", "appreciative", "cheap", "expensive",
		"adventurous", "crude", "fancy", "cruel", "fancy", "omniscient",
	}
	plants = []string{
		"monstera", "philodendron", "pothos", "rex begonia",
		"fig", "mango", "apple", "jalepeno", "habenero",
		"coffee", "snake plant", "spider plant", "curry leaf",
		"achacha", "orange", "lime", "longan", "lychee",
	}
	pref = []string{"low", "medium", "high"}
)

func generatePlants() []Plant {
	pp := []Plant{}
	num := rand.Intn(max-min) + min
	for i := 0; i < num; i++ {
		adj := adjectives[rand.Intn(len(adjectives))]
		plant := plants[rand.Intn(len(plants))]

		pp = append(pp, Plant{
			ID:            uuid.New(),
			Name:          strings.Join([]string{adj, plant}, " "),
			BotanicalName: plant,
			Status:        "Fertilized",
			WaterPref:     pref[rand.Intn(len(pref))],
			LightPref:     pref[rand.Intn(len(pref))],
			HumidityPref:  pref[rand.Intn(len(pref))],
		})
	}
	return pp
}
