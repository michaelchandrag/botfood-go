package http

type HTTPDelivery struct {
	Success		bool 				`json:"success"`
	Data 		interface{} 		`json:"data"`
}