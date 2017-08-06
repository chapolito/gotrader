package main

type Order struct {
	Type string
	Id string
	Price float64
	Size  float64
}

type Orders []Order
