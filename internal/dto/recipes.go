package dto

import "mime/multipart"

type CreateRecipeRequest struct {
	Recipe       recipe        `json:"recipe" form:"recipe"`
	Ingredients  []ingredient  `json:"ingredients"`
	Instructions []instruction `json:"instructions"`
	Image *multipart.FileHeader `form:"image"`
}

type ingredient struct {
	Name            string  `json:"name"`
	Quantity        float64 `json:"quantity"`
	MeasurementUnit string  `json:"measurement_unit"`
}

type instruction struct {
	Content   string `json:"content"`
	Numbering int    `json:"numbering"`
}

type recipe struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CookTime    int    `json:"cook_time"`
	Servings    int    `json:"servings"`
}
