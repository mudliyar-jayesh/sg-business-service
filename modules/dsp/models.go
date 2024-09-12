package dsp

type DSP struct {
	CircleName   string `bson:"CircleName"`
	RegionName   string `bson:"RegionName"`
	DivisionName string `bson:"DivisionName"`
	OfficeName   string `bson:"OfficeName"`
	Pincode      string `bson:"Pincode"`
	OfficeType   string `bson:"OfficeType"`
	Delivery     string `bson:"Delivery"`
	District     string `bson:"District"`
	StateName    string `bson:"StateName"`
	Latitude     string `bson:"Latitude"`
	Longitude    string `bson:"Longitude"`
}
