package main

import "fmt"

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"database"`
}

type DBConfigs struct {
	Source DBConfig `json:"source"`
	Target DBConfig `json:"target"`
}

func (c DBConfig) BuildDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}
