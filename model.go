package main

import (
	"slices"
)

type Config struct {
	Issues []Issue `yaml:"issues"`
}

type Issue struct {
	Name           string `yaml:"name"`
	CreationMonths []int  `yaml:"creation_months"`
}

func (i *Issue) IsCreationMonth(month int) bool {
	return slices.Contains(i.CreationMonths, month)
}
