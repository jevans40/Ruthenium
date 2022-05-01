package constants

const (
	//This will act as a multiple for all communication channels, 0 will disable the channel buffers, and will allow for race testing
	//1 will allow for normal channel buffer size
	RACECHANNELSIZETEST int = 1
)
