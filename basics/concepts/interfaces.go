package concepts

import (
	"fmt"
	"reflect"
	"strconv"
)

// An interface in Golang defines a set of method signatures (the behavior) without implementing them.
// It describes what a type can do, not how it does it.

// -------------------------- interfaces -------------------------

// Base interface for every vehicle
type Vehicle interface {
	TurnOn()
	TurnOff()
	GetPassengerCapacity() string // "Getter"
	GetMaxSpeed() string
	GetBrand() string
	GetColor() string
	GetPrice() string
}

type AerialVehicle interface {
	Vehicle
	TakeOff()
	Land()
	GetMaximunAltitude() int
}

type MilitaryVehicle interface {
	GetAmmunition() int
	LoadAmmunition() // "Setter"
	Shoot()
}

type MilitaryAircraft interface {
	MilitaryVehicle
	AerialVehicle
	LaunchRocket()
}

// -------------------------- implementing interfaces -------------------------
// Golang uses "implicit implementation/interfaces" (duck typing).
// A struct does not need to declare that it implements an interface, it just needs to have matching methods.
// Struct embedding + Interface = Polymorphism

// Base struct for every vehicle
type BaseVehicle struct {
	Brand      string
	Color      string
	Passengers int
	MaxSpeed   int
	Wheels     int
	Price      int
}

func (c BaseVehicle) TurnOn() {
	fmt.Println("Starting the vehicle..")
}

func (BaseVehicle) TurnOff() {
	fmt.Println("Turning off the vehicle..")
}

func (c BaseVehicle) GetPassengerCapacity() string {
	return strconv.Itoa(c.Passengers)
}

func (c BaseVehicle) GetMaxSpeed() string {
	return strconv.Itoa(c.MaxSpeed)
}

func (c BaseVehicle) GetBrand() string {
	return c.Brand
}

func (c BaseVehicle) GetColor() string {
	return c.Color
}

func (c BaseVehicle) GetPrice() string {
	return strconv.Itoa(c.Price)
}

// Base struct for air vehicles
type BaseAerialVehicle struct {
	// Anonymous struct embedded (composition)
	BaseVehicle
	MaximunAltitude int
}

// method shadowing (overriding behavior)
func (c BaseAerialVehicle) TurnOn() {
	fmt.Println("Starting the turbine..")
}

func (c BaseAerialVehicle) TurnOff() {
	fmt.Println("Shutting down the turbine..")
}

func (c BaseAerialVehicle) TakeOff() {
	fmt.Println("Taking off..")
}

func (c BaseAerialVehicle) Land() {
	fmt.Println("Landing..")
}

func (c BaseAerialVehicle) GetMaximunAltitude() int {
	return c.MaximunAltitude
}

type BaseMilitaryVehicle struct {
	Ammunition int
}

func (v BaseMilitaryVehicle) GetAmmunition() int {
	return v.Ammunition
}

func (v *BaseMilitaryVehicle) LoadAmmunition() {
	v.Ammunition++
}

func (v *BaseMilitaryVehicle) Shoot() {
	fmt.Println("Shooting.. -> -> -> ->")
	v.Ammunition--
}

type BaseMilitaryAerialVehicle struct {
	BaseMilitaryVehicle
	BaseAerialVehicle
}

func (v BaseMilitaryAerialVehicle) LaunchRocket() {
	fmt.Println("Rocket launched.. ----------->")
}

// -------------------------- Structs examples -------------------------
type Car struct {
	// Anonymous embedded (composition)
	BaseVehicle
}

// method shadowing (overriding behavior)
func (c Car) TurnOn() {
	fmt.Println("Turning on the engine..")
}

type Bicycle struct {
	BaseVehicle
}

func (b Bicycle) TurnOn() {
	fmt.Println("This bicycle doesn't have a motor")
}

func (b Bicycle) TurnOff() {
	fmt.Println("This bicycle doesn't have a motor")
}

type FireEngine struct {
	BaseVehicle
	WaterCapacity int
	HasLadder     bool
	LadderLength  int
}

type Tank struct {
	BaseVehicle
	BaseMilitaryVehicle
	ArmorLevel int
}

type F22Raptor struct {
	BaseMilitaryAerialVehicle
	StealthCapacity bool
}

// -------------------------- Polymorphism -------------------------
// one interface, multiple implementations
// allows different types of objects to be accessed through a single, shared interface.
// each object can implement that interface in its own way.

func describeVehicle(v Vehicle) {
	structType := reflect.TypeOf(v).Name()
	description := "This vehicle is a " + v.GetColor() + " " + v.GetBrand() + " " + structType
	description += ", capable to reach speeds about " + v.GetMaxSpeed() + " km/h"
	description += ", with capacity for " + v.GetPassengerCapacity() + " passengers"
	description += " and a price tag of " + v.GetPrice()

	fmt.Println(description)
}

func turnOnVehicle(v Vehicle) {
	v.TurnOn()
}

func shoot(v MilitaryVehicle) {
	v.Shoot()
}

func takeOff(v AerialVehicle) {
	v.TakeOff()
}

func launchRocket(v MilitaryAircraft) {
	v.LaunchRocket()
}

// ---------------------------------------------------------------------

func getInstances() (Car, Bicycle, Tank, FireEngine, F22Raptor) {
	car := Car{
		BaseVehicle: BaseVehicle{
			Brand:      "Toyota",
			Color:      "white",
			Passengers: 6,
			MaxSpeed:   240,
			Wheels:     4,
			Price:      80000,
		},
	}

	bike := Bicycle{
		BaseVehicle: BaseVehicle{
			Brand:      "Trek",
			Color:      "gray",
			Passengers: 1,
			MaxSpeed:   120,
			Wheels:     2,
			Price:      2000,
		},
	}

	fireCar := FireEngine{
		BaseVehicle: BaseVehicle{
			Brand:      "Rosenbauer",
			Color:      "read and white",
			Passengers: 12,
			MaxSpeed:   100,
			Wheels:     8,
			Price:      110000,
		},
		WaterCapacity: 20000,
		HasLadder:     true,
		LadderLength:  20,
	}

	tank := Tank{
		ArmorLevel: 7,
		BaseMilitaryVehicle: BaseMilitaryVehicle{
			Ammunition: 0,
		},
		BaseVehicle: BaseVehicle{
			Brand:      "AbramsX",
			Color:      "gray and black",
			Passengers: 4,
			MaxSpeed:   140,
			Wheels:     14,
			Price:      700000,
		},
	}

	aircraft := F22Raptor{
		StealthCapacity: true,
		BaseMilitaryAerialVehicle: BaseMilitaryAerialVehicle{
			BaseMilitaryVehicle: BaseMilitaryVehicle{
				Ammunition: 0,
			},
			BaseAerialVehicle: BaseAerialVehicle{
				MaximunAltitude: 50000,
				BaseVehicle: BaseVehicle{
					Brand:      "Lockheed Martin",
					Color:      "gray",
					Passengers: 1,
					MaxSpeed:   2200,
					Wheels:     4,
					Price:      20000000,
				},
			},
		},
	}

	return car, bike, tank, fireCar, aircraft
}

func InterfacesinAction() {
	car, bike, tank, fireCar, f22Raptor := getInstances()

	describeVehicle(car)
	turnOnVehicle(car)
	// shoot(car) // car cannot fire because it does not "implement" the correct interface

	describeVehicle(bike)
	turnOnVehicle(bike)

	describeVehicle(tank)
	turnOnVehicle(tank)
	shoot(&tank)
	// launchRocket(tank) // tank cannot fire a rocket

	describeVehicle(fireCar)
	turnOnVehicle(fireCar)

	describeVehicle(f22Raptor)
	turnOnVehicle(f22Raptor)
	takeOff(f22Raptor)
	shoot(&f22Raptor)
	launchRocket(&f22Raptor)
}
