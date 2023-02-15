package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Plant struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	BotanicalName string    `json:"botanical_name"`
	Status        string    `json:"status"`
	WaterPref     string    `json:"water_preference"`
	LightPref     string    `json:"light_preference"`
	HumidityPref  string    `json:"humidity_preference"`
}

func (i *Instance) GetAllPlants(ctx context.Context) ([]Plant, error) {
	pp := []Plant{}
	q := fmt.Sprintf(`SELECT 
    id, name, botanical_name, status, water_pref, light_pref, humidity
    FROM %s`, plantsTableName)
	rows, err := i.db.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to get plants: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, name, botanical string
		var status, water, light, humid sql.NullString
		if err = rows.Scan(&id, &name, &botanical, &status, &water, &light, &humid, &humid); err != nil {
			return nil, err
		}
		pp = append(pp, Plant{
			ID:            uuid.MustParse(id),
			Name:          name,
			BotanicalName: botanical,
			Status:        status.String,
			WaterPref:     water.String,
			LightPref:     light.String,
			HumidityPref:  humid.String,
		})
	}

	return pp, nil
}
