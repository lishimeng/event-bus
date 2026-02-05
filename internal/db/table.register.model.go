package db

var CommonModels []any

func init() {

	CommonModels = append(CommonModels,
		new(DataRecord),
		new(ChannelConfig),
		new(SysConfig),
	)
}
