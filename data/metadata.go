package data

import "github.com/spf13/viper"

type MetadataViewResponse struct {
	// use a struct incase we want to add more fields such as 'type'
	Id string `json:"id"`
	MetadataImpl
}

// A placeholder for an implementation of how we deal with metadata
type Metadata interface {
	GetValue(string) interface{}
	SetValue(string, interface{}) Metadata
	GetString(name string) string
	GetInt(name string) int64
	GetBool(name string) bool
	GetFloat(name string) float64
	Merge(Metadata)
	GetMap() map[string]interface{}
}

type MetadataImpl struct {
	Values *viper.Viper `json:"values"`
}

func NewMetadata(baseline ...Metadata) Metadata {
	m := &MetadataImpl{
		Values: viper.New(),
	}

	for _, baselineMetadata := range baseline {
		m.Merge(baselineMetadata)
	}

	return m
}

func (md *MetadataImpl) GetString(name string) string {
	if v, ok := md.GetValue(name).(string); ok {
		return v
	}
	return ""
}

func (md *MetadataImpl) GetValue(name string) interface{} {
	return md.Values.Get(name)
}

func (md *MetadataImpl) GetInt(name string) int64 {
	return md.Values.GetInt64(name)
}

func (md *MetadataImpl) GetBool(name string) bool {
	return md.Values.GetBool(name)
}

func (md *MetadataImpl) GetFloat(name string) float64 {
	return md.Values.GetFloat64(name)
}

func (md *MetadataImpl) SetValue(name string, value interface{}) Metadata {
	md.Values.Set(name, value)
	return md
}

func (md *MetadataImpl) Merge(other Metadata) {
	md.Values.MergeConfigMap(other.(*MetadataImpl).Values.AllSettings())
}

func (md *MetadataImpl) GetMap() map[string]interface{} {
	return md.Values.AllSettings()
}
