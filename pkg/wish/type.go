package wish

type Type struct {
	Key  int    `json:"key,string" mapstructure:"key"`
	Name string `json:"name" mapstructure:"name"`
}
