package domain

// Car representa la información del vehículo utilizado en el viaje
type Car struct {
	Brand string `json:"brand" bson:"brand" binding:"required"`
	Model string `json:"model" bson:"model" binding:"required"`
	Year  int    `json:"year" bson:"year" binding:"required"`
	Color string `json:"color" bson:"color" binding:"required"`
	Plate string `json:"plate" bson:"plate" binding:"required"`
}
