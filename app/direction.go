package app

type Direction int8

const (
	DIn   Direction = -1
	DNegY Direction = 0 // (0,-1,0)
	DPosY Direction = 1 // (0,1,0)
	DNegZ Direction = 2 // (0,0,-1)
	DPosZ Direction = 3 // (0,0,1)
	DPosX Direction = 4 // (-1,0,0)
	DNegX Direction = 5 // (1,0,0)
)
