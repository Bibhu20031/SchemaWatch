package schema

type RegisterRequest struct {
	DBHost     string `json:"db_host" binding:"required"`
	DBPort     int    `json:"db_port" binding:"required"`
	DBName     string `json:"db_name" binding:"required"`
	SchemaName string `json:"schema_name" binding:"required"`
	TableName  string `json:"table_name" binding:"required"`
}
